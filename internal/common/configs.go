package common

import "gopkg.in/yaml.v3"

type DatabaseConfig struct {
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	Name     string        `yaml:"name"`
	Schema   string        `yaml:"schema"`
	User     string        `yaml:"user"`
	Password string        `yaml:"password"`
	Seed     int64         `yaml:"seed"`
	Tables   []TableConfig `yaml:"tables"`
}

func (dc *DatabaseConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, dc)
}

type TableConfig struct {
	Name    string         `yaml:"name"`
	Columns []ColumnConfig `yaml:"columns"`
}

func (tc *TableConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, tc)
}

type ColumnConfig struct {
	Name            string            `yaml:"name"`
	Type            PrimitiveDataType `yaml:"type"`
	Nullable        bool              `yaml:"nullable"`
	MaxLength       int               `yaml:"max-length,omitempty"`
	DefaultValue    string            `yaml:"default-value,omitempty"`
	BlankPercentage int               `yaml:"blank-percentage,omitempty"`
	GeneratorType   DataGeneratorType `yaml:"generator"`
	Options         map[string]any    `yaml:"options,omitempty"`
}

func (cc *ColumnConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, cc)
}
