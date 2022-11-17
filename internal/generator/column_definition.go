package generator

import (
	"dbgo/internal/common"
	"fmt"
	"strings"
)

type ColumnDefinition struct {
	Name            string
	Type            common.PrimitiveDataType
	Nullable        bool
	MaxLength       int
	DefaultValue    any
	BlankPercentage int
	Generator       DataGenerator
	Options         map[string]any
}

func (c *ColumnDefinition) OptionAsString(key string, defaultValue string) string {
	if v, found := c.Options[key]; found {
		return v.(string)
	}
	return defaultValue
}

func (c *ColumnDefinition) OptionAsInt(key string, defaultValue int) int {
	if v, found := c.Options[key]; found {
		return v.(int)
	}
	return defaultValue
}

func ProcessColumns(columnConfig []common.ColumnConfig) []ColumnDefinition {
	var filteredColumns []ColumnDefinition
	for columnIdx := 0; columnIdx < len(columnConfig); columnIdx++ {
		cfg := columnConfig[columnIdx]
		if !cfg.GeneratorType.FilterParam() {
			filteredColumns = append(filteredColumns, ColumnDefinition{
				Name:            cfg.Name,
				Type:            cfg.Type,
				Nullable:        cfg.Nullable,
				MaxLength:       cfg.MaxLength,
				DefaultValue:    cfg.DefaultValue,
				BlankPercentage: cfg.BlankPercentage,
				Generator:       CreateGeneratorFromType(cfg.GeneratorType, cfg.Options), // todo, should probably be a sub map of generator?
				Options:         cfg.Options,
			})
		}
	}
	return filteredColumns
}

func ExtractColumnNames(columns []ColumnDefinition) string {
	var names []string
	for i := 0; i < len(columns); i++ {
		names = append(names, columns[i].Name)
	}

	return fmt.Sprintf("%s", strings.Join(names, ","))
}

func ExtractValuePlaceholders(columns []ColumnDefinition) string {
	var placeholders []string
	for i := 0; i < len(columns); i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	return fmt.Sprintf("%s", strings.Join(placeholders, ","))
}

func ExtractColumnValues(columns []ColumnDefinition, rowNum int) []any {
	var values []any
	for i := 0; i < len(columns); i++ {
		values = append(values, columns[i].Generator.Generate(columns[i], rowNum))
	}

	return values
}
