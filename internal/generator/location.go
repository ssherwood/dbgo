package generator

import (
	"dbgo/internal/common"
	"dbgo/pkg/randoms"
)

type LocationGenerator struct {
	Options map[string]any
}

func (g *LocationGenerator) Generate(column ColumnDefinition, rowNumber int) any {
	generatorType := column.OptionAsString(common.OptionsGeneratorType, common.LocationGeneratorTypeCity)

	switch generatorType {
	case common.LocationGeneratorTypeCity:
		return randoms.RandomCity("US-NY").Name
	}

	return column.DefaultValue
}
