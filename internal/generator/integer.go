package generator

import (
	"dbgo/internal/common"
	"dbgo/pkg/randoms"
)

type IntegerGenerator struct {
	Options map[string]any
}

func (g *IntegerGenerator) Generate(column ColumnDefinition, rowNumber int) any {
	generatorType := column.OptionAsString(common.OptionsGeneratorType, common.IntegerGeneratorTypeRandom)
	maxValue := column.OptionAsInt(common.IntegerGeneratorOptionMaxValue, 32767)
	minValue := column.OptionAsInt(common.IntegerGeneratorOptionMinValue, 0)

	if minValue > maxValue {
		minValue = 0 // TODO warn
	}

	switch generatorType {
	case common.IntegerGeneratorTypeRandom:
		return randoms.RandomIntBetween(maxValue, minValue)
	case common.IntegerGeneratorTypeSequential:
		return (rowNumber % (maxValue - minValue)) + minValue
	}

	return column.DefaultValue
}
