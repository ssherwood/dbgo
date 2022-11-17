package generator

type DefaultGenerator struct {
	Options map[string]any
}

func (g *DefaultGenerator) Generate(column ColumnDefinition, rowNumber int) any {
	return column.Type.DefaultValue(column.MaxLength)
}
