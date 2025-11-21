package zlog

import (
	"fmt"
)

type LoggerConfig struct {
	Level      Level             `yaml:"level"`
	Output     string            `yaml:"output"` // file、console、both
	Format     string            `yaml:"format"` // json、console
	FilePath   string            `yaml:"file_path"`
	MaxSize    int               `yaml:"max_size"`
	MaxBackups int               `yaml:"max_backups"`
	MaxAge     int               `yaml:"max_age"`
	Compress   bool              `yaml:"compress"`
	Sampling   bool              `yaml:"sampling"`
	Fields     map[string]string `yaml:"fields"`
}

func (c *LoggerConfig) Validate() error {
	if c.MaxSize <= 0 {
		c.MaxSize = 100
	}
	if c.MaxBackups < 0 {
		c.MaxBackups = 10
	}
	if c.MaxAge < 0 {
		c.MaxAge = 30
	}
	if (c.Output == "file" || c.Output == "both") && c.FilePath == "" {
		return fmt.Errorf("FilePath is required when Output='file'")
	}
	return nil
}

func DefaultConfig() LoggerConfig {
	return LoggerConfig{
		Level:      InfoLevel,
		Output:     "console",
		Format:     "console",
		FilePath:   "",
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
		Sampling:   false,
		Fields:     make(map[string]string),
	}
}
