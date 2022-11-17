package generator

import (
	"dbgo/pkg/randoms"
	"time"
)

type TimestampGenerator struct {
	Options map[string]any
}

func (g *TimestampGenerator) Generate(column ColumnDefinition, rowNumber int) any {
	// TODO allow config to override default format of RFC3339
	return randoms.RandomDateUTC(time.Now().Year(), 1990).Format(time.RFC3339)
}
