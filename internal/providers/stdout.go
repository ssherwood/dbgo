package providers

import (
	"dbgo/internal/common"
	"dbgo/internal/generator"
	"fmt"
	"io"
	"strings"
)

const (
	stdoutInsertSingleRow = "INSERT INTO %s.%s(%s) VALUES(%s);\n"
	stdoutBeginTx         = "BEGIN;\n"
	stdoutEndTx           = "COMMIT;\n"
)

type StdoutDBProvider struct {
	Writer io.Writer
}

func (p *StdoutDBProvider) Init() error {
	return nil
}

func (p *StdoutDBProvider) Close() {
}

func (p *StdoutDBProvider) Discover(schema string) ([]common.TableConfig, error) {
	return nil, nil // not implementable
}

func (p *StdoutDBProvider) InsertSingleRow(schema string, table string, columns []generator.ColumnDefinition, iteration int, transactional bool) error {
	statement := fmt.Sprintf(stdoutInsertSingleRow, schema, table, generator.ExtractColumnNames(columns), strings.Join(extractColValues(columns, iteration), ","))
	if transactional {
		p.Writer.Write([]byte(stdoutBeginTx))
		p.Writer.Write([]byte(statement))
		p.Writer.Write([]byte(stdoutEndTx))
	} else {
		p.Writer.Write([]byte(statement))
	}

	return nil
}

func (p *StdoutDBProvider) InsertBatchRow(schema string, table string, columns []generator.ColumnDefinition, iteration int, batchSize int) error {
	//TODO implement me
	panic("implement me")
}

func extractColValues(columns []generator.ColumnDefinition, rowNum int) []string {
	var values []string
	for i := 0; i < len(columns); i++ {
		columnValue := fmt.Sprintf("%v", columns[i].Generator.Generate(columns[i], rowNum))
		if columns[i].Type.QuotedSQL() {
			values = append(values, fmt.Sprintf("'%s'", strings.ReplaceAll(columnValue, "'", `''`)))
		} else {
			if columnValue == "<nil>" {
				values = append(values, "NULL")
			} else {
				values = append(values, columnValue)
			}
		}
	}
	return values
}
