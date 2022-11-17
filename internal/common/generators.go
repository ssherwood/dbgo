package common

type DataGeneratorType string

const (
	GeneratorTypeDefault    DataGeneratorType = "default"
	GeneratorTypeSerial     DataGeneratorType = "serial"
	GeneratorTypeSequence   DataGeneratorType = "sequence"
	GeneratorTypeInteger    DataGeneratorType = "integer"
	GeneratorTypeDecimal    DataGeneratorType = "decimal"
	GeneratorTypeBoolean    DataGeneratorType = "boolean"
	GeneratorTypeString     DataGeneratorType = "string"
	GeneratorTypeCustomList DataGeneratorType = "custom-list"
	GeneratorTypeTimestamp  DataGeneratorType = "timestamp"
	GeneratorTypeUUID       DataGeneratorType = "uuid"
)

const (
	OptionsGeneratorType = "generator-type"

	SerialGeneratorOptionStartValue = "start-value"

	SequenceGeneratorOptionName = "sequence-name"

	IntegerGeneratorTypeRandom     = "random"
	IntegerGeneratorTypeSequential = "sequential"
	IntegerGeneratorOptionMinValue = "min-value"
	IntegerGeneratorOptionMaxValue = "max-value"

	StringGeneratorTypeAscii       = "ascii"
	StringGeneratorTypeUtf8        = "utf8"
	StringGeneratorOptionMaxLength = "max-length"
	StringGeneratorOptionMinLength = "min-length"

	CustomListTypeSequential = "sequential"
	CustomListTypeWeighted   = "weighted"
	CustomListTypeRandom     = "random"
	CustomListOptionValues   = "values"
)

//
// FilterParam returns true if the generator type determines if it should be
// excluded from the query params (this is typically only for the "default"
// generator which implies it will be handled using the native database
// default).
//
func (g *DataGeneratorType) FilterParam() bool {
	return *g == GeneratorTypeDefault
}

//
// DefaultOptions produces a basic option defaults for a given generator type.
// These have very limited visibility into the data types or other options that
// may be further used to customize the generator behavior.
//
// Note, maxSize may represent the maximum length or maximum bit size depending
// on the generator type.
//
func (g *DataGeneratorType) DefaultOptions(maxSize int) map[string]any {
	switch *g {
	case GeneratorTypeString:
		return map[string]any{
			OptionsGeneratorType:           StringGeneratorTypeAscii,
			StringGeneratorOptionMaxLength: maxSize,
			StringGeneratorOptionMinLength: 1,
		}
	case GeneratorTypeInteger:
		return map[string]any{
			OptionsGeneratorType:           IntegerGeneratorTypeRandom,
			IntegerGeneratorOptionMaxValue: calculateMaxValueFromBits(maxSize),
			IntegerGeneratorOptionMinValue: 0,
		}
	}

	return map[string]any{}
}

func calculateMaxValueFromBits(maxSize int) int {
	maxValue := 1
	if maxSize == 16 {
		maxValue = 32767
	} else if maxSize == 32 {
		maxValue = 2147483647
	} else if maxSize == 64 {
		maxSize = 9223372036854775807
	}
	return maxValue
}
