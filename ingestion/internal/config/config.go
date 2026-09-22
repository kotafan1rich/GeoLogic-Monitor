package config

import (
	"fmt"
	"time"
)

type Config struct {
	Logger     LoggerConfig     `yaml:"logger"`
	Aggregator AggregatorConfig `yaml:"aggregator"`
	GeoApi     GeoApiConfig     `yaml:"geo-api"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type AggregatorConfig struct {
	DigitalSpb DigitalSpbConfig `yaml:"digitalspb"`
}

type DigitalSpbConfig struct {
	MaxIdleConns        int               `yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int               `yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost     int               `yaml:"max_conns_per_host"`
	RequestTimeout      time.Duration     `yaml:"request_timeout"`
	AttemptTimeout      time.Duration     `yaml:"attempt_timeout"`
	MaxRetries          uint              `yaml:"max_retries"`
	JobInterval         time.Duration     `yaml:"job_interval"`
	BaseURLMap          map[string]string `yaml:"base_urls"`
	StaticFiles         map[string]string `yaml:"static_files"`
}

type GeoApiConfig struct {
	URL                 string        `yaml:"url"`
	AuthToken           string        `env:"INGESTION_SERVICE_TOKEN" env-required:"true"`
	MaxIdleConns        int           `yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int           `yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost     int           `yaml:"max_conns_per_host"`
	RequestTimeout      time.Duration `yaml:"request_timeout"`
	AttemptTimeout      time.Duration `yaml:"attempt_timeout"`
	MaxRetries          uint          `yaml:"max_retries"`
	InfraTypes          []InfraType   `yaml:"infra_types"`
}

type InfraType struct {
	Slug      string `yaml:"slug"`
	Name      string `yaml:"name"`
	Weight    int    `yaml:"weight"`
	MaxRadius int    `yaml:"max_radius"`
}

var config *Config

func Get() *Config {
	const op = "ingestion.config.Get"

	if config == nil {
		panic(fmt.Sprintf("%s: config is not initialized", op))
	}
	return config
}
