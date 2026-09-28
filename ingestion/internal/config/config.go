package config

import (
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	Logger     LoggerConfig
	Database   DatabaseConfig `yaml:"database"`
	Producer   ProducerConfig `yaml:"producer"`
	Scheduler  SchedulerConfig
	Aggregator AggregatorConfig `yaml:"aggregator"`
	GeoApi     GeoApiConfig     `yaml:"geo-api"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

type SchedulerConfig struct {
	Infra              string `env:"SCHEDULER_INFRA_CRON"       env-default:"0 3 1 * *"`
	Events             string `env:"SCHEDULER_EVENTS_CRON"      env-default:"0 4 * * 1"`
	Competitors        string `env:"SCHEDULER_COMPETITORS_CRON" env-default:"0 6 * * *"`
	EventNotifications string `env:"SCHEDULER_EVENT_NOTIFICATIONS_CRON" env-default:"0 7 * * *"`
}

type LoggerConfig struct {
	Level  string `env:"LOG_LEVEL"  env-default:"info"`
	Format string `env:"LOG_FORMAT" env-default:"json"`
}

type DatabaseConfig struct {
	Postgresql PostgresqlConfig `yaml:"postgresql"`
}

type ProducerConfig struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

type PostgresqlConfig struct {
	Host                string        `env:"POSTGRES_HOST" env-required:"true"`
	Port                int           `env:"POSTGRES_PORT" env-required:"true"`
	User                string        `env:"POSTGRES_USER" env-required:"true"`
	Password            string        `env:"POSTGRES_PASSWORD" env-required:"true"`
	DB                  string        `env:"POSTGRES_DB" env-required:"true"`
	SSLMode             string        `yaml:"ssl_mode"`
	MinConns            int           `yaml:"min_conns"`
	MaxConns            int           `yaml:"max_conns"`
	MaxConnIdleLifetime time.Duration `yaml:"max_conn_idle_lifetime"`
	MaxConnLifetime     time.Duration `yaml:"max_conn_lifetime"`
}

type KafkaConfig struct {
	Addresses      []string      `env:"KAFKA_ADDRESSES" env-required:"true"`
	Topic          string        `yaml:"topic" env-default:"notifications"`
	FlushTimeout   time.Duration `yaml:"flush_timeout"`
	AttemptTimeout time.Duration `yaml:"attempt_timeout"`
	MaxRetries     uint          `yaml:"max_retries"`
	BatchSize      int32         `yaml:"batch_size"`
	BufferMaxMsg   int           `yaml:"buffer_max_msg"`
	BufferMaxBytes int           `yaml:"buffer_max_bytes"`
}

func (c *PostgresqlConfig) DSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, strconv.Itoa(c.Port), c.DB, c.SSLMode,
	)
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
	Maps       MapsConfig       `yaml:"maps"`
	TwoGis     TwoGisConfig     `yaml:"twogis"`
}

type DigitalSpbConfig struct {
	HTTP        HTTPClientConfig  `yaml:"http"`
	BaseURLMap  map[string]string `yaml:"base_urls"`
	StaticFiles map[string]string `yaml:"static_files"`
}

type MapsConfig struct {
	HTTP    HTTPClientConfig `yaml:"http"`
	BaseURL string           `yaml:"base_url"`
}

type TwoGisConfig struct {
	HTTP    HTTPClientConfig `yaml:"http"`
	BaseURL string           `yaml:"base_url" env-default:"https://catalog.api.2gis.com/"`
	APIKey  string           `env:"TWOGIS_API_KEY" env-required:"true"`
}

type MonitoringConfig struct {
	FetchConcurrency int                         `yaml:"fetch_concurrency" env-default:"2"`
	RouteConcurrency int                         `yaml:"route_concurrency" env-default:"4"`
	Competitors      CompetitorsMonitoringConfig `yaml:"competitors"`
	Events           EventsMonitoringConfig      `yaml:"events"`
}

type CompetitorsMonitoringConfig struct {
	Window              time.Duration `yaml:"window"               env-default:"2160h"`
	CheckpointBootstrap time.Duration `yaml:"checkpoint_bootstrap" env-default:"24h"`
	TTL                 time.Duration `yaml:"ttl"                  env-default:"24h"`
}

type EventsMonitoringConfig struct {
	Window       time.Duration `yaml:"window"        env-default:"24h"`
	SearchRadius int           `yaml:"search_radius" env-default:"1200"`
	WalkRadius   int           `yaml:"walk_radius"   env-default:"800"`
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
	Slug      string   `yaml:"slug"`
	Name      string   `yaml:"name"`
	Weight    int      `yaml:"weight"`
	MaxRadius int      `yaml:"max_radius"`
	Query     string   `yaml:"query"`
	Rubrics   []string `yaml:"rubrics"`
}

var config *Config

func Get() *Config {
	const op = "ingestion.config.Get"

	if config == nil {
		panic(fmt.Sprintf("%s: config is not initialized", op))
	}
	return config
}
