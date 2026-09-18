package config

import (
	"fmt"
	"time"
)

type Config struct {
	Logger     LoggerConfig     `yaml:"logger"`
	Aggregator AggregatorConfig `yaml:"aggregator"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type AggregatorConfig struct {
	DigitalSpb DigitalSpbConfig `yaml:"digitalspb"`
}

type DigitalSpbConfig struct {
	MaxIdleConns        int           `yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int           `yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost     int           `yaml:"max_conns_per_host"`
	RequestTimeout      time.Duration `yaml:"request_timeout"`
}

var config *Config

func Get() *Config {
	const op = "ingestion.config.Get"

	if config == nil {
		panic(fmt.Sprintf("%s: config is not initialized", op))
	}
	return config
}
