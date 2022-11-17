package providers

import (
	"dbgo/internal/common"
	"dbgo/internal/generator"
)

type DatabaseProvider interface {
	Init() error
	Close()
	Discover(schema string) ([]common.TableConfig, error)
	InsertSingleRow(schema string, table string, columns []generator.ColumnDefinition, iteration int, transactional bool) error
	InsertBatchRow(schema string, table string, columns []generator.ColumnDefinition, iteration int, batchSize int) error
}
