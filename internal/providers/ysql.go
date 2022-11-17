package providers

import (
	"context"
	"dbgo/internal/common"
	"dbgo/internal/generator"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/exp/maps"
	"strings"
)

const (
	ysqlInsertSingleRow = "INSERT INTO %s.%s(%s) VALUES(%s);"
)

type YugabyteYSQL struct {
	Config common.DatabaseConfig
	pool   *pgxpool.Pool
}

func (yb *YugabyteYSQL) Init() error {
	//statement_cache_capacity=1
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", yb.Config.User, yb.Config.Password, yb.Config.Host, yb.Config.Port, yb.Config.Name)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}

	// _0: INSERT INTO public.users(id,account_status,city,state,auto_renew,about,encryption_level) VALUES│
	//config.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
	//	_, err := c.Prepare(ctx, "ps1", "INSERT INTO public.users(id,account_status,city,state,auto_renew,about,encryption_level) VALUES($1,$2,$3,$4,$5,$6,$7);")
	//	return err
	//}

	//logger := zap.NewExample()
	//defer logger.Sync()
	//config.ConnConfig.Logger = zapadapter.NewLogger(logger)

	//fmt.Println(config)

	yb.pool, err = pgxpool.ConnectConfig(context.Background(), config)
	return err
}

func (yb *YugabyteYSQL) Close() {
	yb.pool.Close()
}

//
// InsertSingleRow
// Simulates a single row auto-commit style insert
//
func (yb *YugabyteYSQL) InsertSingleRow(schema string, table string, columns []generator.ColumnDefinition, iteration int, transactional bool) error {
	statement := fmt.Sprintf(ysqlInsertSingleRow, schema, table, generator.ExtractColumnNames(columns), generator.ExtractValuePlaceholders(columns))
	values := generator.ExtractColumnValues(columns, iteration)

	conn, err := yb.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	// todo test client-side timeout
	//conn.Conn().PgConn().Conn().SetDeadline(time.Now().Add(time.Millisecond * 500))

	// TODO does pgx actually prepare and cache statements for us???
	// docs say they do it automatically but not sure

	if transactional {
		tx, errTx := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if errTx != nil {
			return errTx
		}
		defer func() {
			if errTx != nil {
				tx.Rollback(context.Background())
			} else {
				tx.Commit(context.Background())
			}
		}()

		_, errTx = tx.Exec(context.Background(), statement, values...)
		return errTx
	} else {
		//_, err = conn.Conn().PgConn().ExecPrepared(context.Background(), "ps1", values..., ni)
		_, err = conn.Exec(context.Background(), statement, values...)
		return err
	}
}

//
// InsertBatchRow
// This isn't performing as well as single inserts, may need to look into this...
//
func (yb *YugabyteYSQL) InsertBatchRow(schema string, table string, columns []generator.ColumnDefinition, iteration int, batchSize int) error {
	batch := &pgx.Batch{}

	for i := 0; i < batchSize; i++ {
		statement := fmt.Sprintf(ysqlInsertSingleRow, schema, table, generator.ExtractColumnNames(columns), generator.ExtractValuePlaceholders(columns))
		values := generator.ExtractColumnValues(columns, iteration+i)
		batch.Queue(statement, values...)
	}

	conn, err := yb.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	batchResults := conn.SendBatch(context.Background(), batch)
	_, err = batchResults.Exec()
	return err
}

func (yb *YugabyteYSQL) Discover(schema string) ([]common.TableConfig, error) {

	statement := `select table_name,
					ordinal_position as position,
					column_name,
					data_type,
					case when character_maximum_length is not null then character_maximum_length else numeric_precision end as max_length,
					is_nullable::Boolean,
					coalesce(column_default, '') as default_value
					from information_schema.columns
					where table_schema = $1
					order by table_schema, table_name, ordinal_position;`

	conn, err := yb.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), statement, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnsMetadata := map[string][]PostgresColumnMetadata{}
	for rows.Next() {
		var metadata PostgresColumnMetadata
		if err := rows.Scan(&metadata.TableName, &metadata.Position, &metadata.ColumnName, &metadata.DataType, &metadata.MaxLength, &metadata.IsNullable, &metadata.DefaultValue); err != nil {
			return nil, err
		}

		columnsMetadata[metadata.TableName] = append(columnsMetadata[metadata.TableName], metadata)
	}

	var tableConfigs []common.TableConfig

	for table, metadata := range columnsMetadata {
		var columnConfigs []common.ColumnConfig
		for i := 0; i < len(metadata); i++ {
			maxLength := 0
			if metadata[i].MaxLength != nil {
				maxLength = *metadata[i].MaxLength
			}

			generatorType, customOptions := metadata[i].asGeneratorType()
			maps.Copy(customOptions, generatorType.DefaultOptions(maxLength))

			columnConfigs = append(columnConfigs, common.ColumnConfig{
				Name:            metadata[i].ColumnName,
				Type:            metadata[i].asPrimitiveDataType(),
				Nullable:        metadata[i].IsNullable,
				MaxLength:       maxLength,
				DefaultValue:    metadata[i].asDefaultValue(),
				BlankPercentage: 0,
				GeneratorType:   generatorType,
				Options:         customOptions,
			})
		}

		tableConfigs = append(tableConfigs, common.TableConfig{
			Name:    table,
			Columns: columnConfigs,
		})

	}

	return tableConfigs, nil
}

type PostgresColumnMetadata struct {
	TableName    string
	Position     int
	ColumnName   string
	DataType     string
	MaxLength    *int
	IsNullable   bool
	DefaultValue string
}

var typeMap = map[string]common.PrimitiveDataType{
	"smallint": common.PrimitiveNumeric, // 2 bytes; -32768 to +32767
	"integer":  common.PrimitiveNumeric, // 4 bytes; -2147483648 to +2147483647
	"bigint":   common.PrimitiveNumeric, // 8 bytes; -9223372036854775808 to +9223372036854775807

	"smallserial": common.PrimitiveNumeric, // 2 bytes; 1 to 32767
	"serial":      common.PrimitiveNumeric, // 4 bytes; 1 to +2147483647
	"bigserial":   common.PrimitiveNumeric, // 8 bytes; 1 to +9223372036854775807

	"decimal":          common.PrimitiveDecimal, // variable bytes; up to 16383 decimal digits precision
	"numeric":          common.PrimitiveDecimal, // variable bytes; up to 16383 decimal digits precision
	"real":             common.PrimitiveDecimal, // 4 bytes; 6 decimal digits precision
	"double precision": common.PrimitiveDecimal, // 8 bytes; 15 decimal digits precision

	"boolean": common.PrimitiveBoolean,

	"character":         common.PrimitiveCharacter,
	"character varying": common.PrimitiveCharacter,
	"text":              common.PrimitiveCharacter,

	"timestamp with time zone": common.PrimitiveDate,
}

func (c *PostgresColumnMetadata) asPrimitiveDataType() common.PrimitiveDataType {
	if v, ok := typeMap[c.DataType]; ok {
		return v
	}

	return "unmapped-primitive"
}

func (c *PostgresColumnMetadata) asDefaultValue() string {
	// TODO detect functions, etc

	if strings.Contains(c.DefaultValue, "nextval(") {
		return ""
	}

	return c.DefaultValue
}

//
// asGeneratorType returns the concrete implementation of the generator for a
// given type using the column metadata.  Using the additional metadata allows
// for a greater interpretation of the generator to use.
//
func (c *PostgresColumnMetadata) asGeneratorType() (common.DataGeneratorType, map[string]any) {

	// if the default has a nextval(...) function use a sequence generator
	// based on the name within the function call
	if strings.Contains(c.DefaultValue, "nextval(") {
		return common.GeneratorTypeSequence, map[string]any{
			common.SequenceGeneratorOptionName: extractSequenceName(c.DefaultValue),
		}
	}

	if c.DataType == "bigint" || c.DataType == "integer" || c.DataType == "smallint" {

		// simple guess that ids without nextval functions will be iterators
		if c.ColumnName == "id" {
			return common.GeneratorTypeSerial, map[string]any{}
		}

		return common.GeneratorTypeInteger, map[string]any{}
	}

	if c.DataType == "real" || c.DataType == "double precision" {
		return common.GeneratorTypeDecimal, map[string]any{}
	}

	if c.DataType == "character varying" {
		//if strings.Contains(c.ColumnName, "city") {
		//	return generator.City
		//}
		//
		//if strings.Contains(c.ColumnName, "state") {
		//	return generator.State
		//}

		return common.GeneratorTypeString, map[string]any{}
	}

	if c.DataType == "boolean" {
		return common.GeneratorTypeBoolean, map[string]any{}
	}

	if c.DataType == "timestamp with time zone" {
		return common.GeneratorTypeTimestamp, map[string]any{}
	}

	if c.DataType == "uuid" {
		return common.GeneratorTypeUUID, map[string]any{}
	}

	return common.GeneratorTypeDefault, map[string]any{} // probably should be blank or null?
}

func extractSequenceName(value string) string {
	strSplit := strings.Split(value, "'") // split between quotes
	sequenceName := strSplit[1]
	return sequenceName
}
