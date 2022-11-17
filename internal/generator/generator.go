package generator

import (
	"dbgo/internal/common"
	"fmt"
)

type DataGenerator interface {
	Generate(column ColumnDefinition, rowNumber int) any
}

func CreateGeneratorFromType(generatorType common.DataGeneratorType, options map[string]any) DataGenerator {
	switch generatorType {
	case common.GeneratorTypeCustomList:
		return &CustomListGenerator{options}
	case common.GeneratorTypeDecimal:
		return &DecimalGenerator{options}
	case common.GeneratorTypeInteger:
		return &IntegerGenerator{options}
	case common.GeneratorTypeSerial:
		return &SerialGenerator{options}
	case common.GeneratorTypeString:
		return &StringGenerator{options}
	case common.GeneratorTypeTimestamp:
		return &TimestampGenerator{options}
	default:
		fmt.Println("!!! Missing generatorType ", generatorType)
		return &DefaultGenerator{options}
	}
}
