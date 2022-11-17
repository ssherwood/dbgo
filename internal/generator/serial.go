package generator

import "dbgo/internal/common"

type SerialGenerator struct {
	Options map[string]any
}

func (g *SerialGenerator) Generate(column ColumnDefinition, rowNumber int) any {
	startValue := column.OptionAsInt(common.SerialGeneratorOptionStartValue, 1)
	return startValue + rowNumber
}
