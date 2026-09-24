package config

import (
	"fmt"
	"time"
)

type Config struct {
	Logger     LoggerConfig
	Scheduler  SchedulerConfig
	Aggregator AggregatorConfig `yaml:"aggregator"`
	GeoApi     GeoApiConfig     `yaml:"geo-api"`
}

type SchedulerConfig struct {
	Infra  string `env:"SCHEDULER_INFRA_CRON"  env-default:"0 3 1 * *"`
	Events string `env:"SCHEDULER_EVENTS_CRON" env-default:"0 4 * * 1"`
}

type LoggerConfig struct {
	Level  string `env:"LOG_LEVEL"  env-default:"info"`
	Format string `env:"LOG_FORMAT" env-default:"json"`
}

type RateLimitConfig struct {
	RPS   float64 `yaml:"rps"`
	Burst int     `yaml:"burst"`
}

type HTTPClientConfig struct {
	MaxIdleConns        int             `yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int             `yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost     int             `yaml:"max_conns_per_host"`
	RequestTimeout      time.Duration   `yaml:"request_timeout"`
	AttemptTimeout      time.Duration   `yaml:"attempt_timeout"`
	MaxRetries          uint            `yaml:"max_retries"`
	RateLimit           RateLimitConfig `yaml:"rate_limit"`
}

type AggregatorConfig struct {
	DigitalSpb DigitalSpbConfig `yaml:"digitalspb"`
}

type DigitalSpbConfig struct {
	HTTP        HTTPClientConfig  `yaml:"http"`
	BaseURLMap  map[string]string `yaml:"base_urls"`
	StaticFiles map[string]string `yaml:"static_files"`
}

type GeoApiConfig struct {
	URL              string           `env:"GEO_API_URL" env-default:"http://api:8080/"`
	AuthToken        string           `env:"INGESTION_SERVICE_TOKEN" env-required:"true"`
	Write            HTTPClientConfig `yaml:"write"`
	Geocoding        HTTPClientConfig `yaml:"geocoding"`
	WriteConcurrency int              `yaml:"write_concurrency"`
	InfraTypes       []InfraType      `yaml:"infra_types"`
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
