package config

import "fmt"

type Config struct {
	Logger LoggerConfig `yaml:"logger"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

var config *Config

func Get() *Config {
	const op = "ingestion.config.Get"

	if config == nil {
		panic(fmt.Sprintf("%s: config is not initialized", op))
	}
	return config
}
