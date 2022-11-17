package generator

type DecimalGenerator struct {
	Options map[string]any
}

func (g *DecimalGenerator) Generate(column ColumnDefinition, rowNumber int) any {
	return column.Type.DefaultValue(column.MaxLength)
}
