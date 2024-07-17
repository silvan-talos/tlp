package config

type Config struct {
	Log         LogConfig         `yaml:"log"`
	Transaction TransactionConfig `yaml:"transaction"`
}

type LogConfig struct {
	Level               string              `yaml:"level"`
	ProcessingType      string              `yaml:"processing"`
	OutputFile          string              `yaml:"output_file"`
	PermanentAttributes []map[string]string `yaml:"permanent_attributes"`
}

type TransactionConfig struct {
	RecorderType string `yaml:"recorder"`
}
