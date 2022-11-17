package generator

import (
	"dbgo/internal/common"
	"dbgo/pkg/randoms"
)

type StringGenerator struct {
	Options map[string]any
}

func (g *StringGenerator) Generate(column ColumnDefinition, rowNum int) any {
	generatorType := column.OptionAsString(common.OptionsGeneratorType, common.StringGeneratorTypeAscii)
	maxLength := column.OptionAsInt(common.StringGeneratorOptionMaxLength, column.MaxLength)
	minLength := column.OptionAsInt(common.StringGeneratorOptionMinLength, maxLength)

	switch generatorType {
	case common.StringGeneratorTypeAscii:
		return randoms.RandomStringAscii(maxLength, minLength)
	case common.StringGeneratorTypeUtf8:
		// todo
	}

	// fall back to the safest string representation
	return randoms.RandomStringAlpha(maxLength)
}
