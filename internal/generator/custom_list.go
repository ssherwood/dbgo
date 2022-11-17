package generator

import (
	"dbgo/internal/common"
	"math/rand"
)

type CustomListGenerator struct {
	Options map[string]any
}

func (c *CustomListGenerator) Generate(column ColumnDefinition, rowNumber int) any {
	generatorType := column.OptionAsString(common.OptionsGeneratorType, common.CustomListTypeRandom)
	listOfValues := column.Options[common.CustomListOptionValues].([]interface{})

	switch generatorType {
	case common.CustomListTypeRandom:
		return listOfValues[rand.Intn(len(listOfValues))]
	case common.CustomListTypeSequential:
		return listOfValues[rowNumber%len(listOfValues)]
	case common.CustomListTypeWeighted:
		// TODO
	}

	return listOfValues[rand.Intn(len(listOfValues))]
}
