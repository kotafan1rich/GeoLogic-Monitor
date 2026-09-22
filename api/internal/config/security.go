package config

import "time"

type Security struct {
	MaxBotToken           string        `env:"MAX_BOT_TOKEN,required"`
	MiniAppInitDataMaxAge time.Duration `env:"MINIAPP_INIT_DATA_MAX_AGE" envDefault:"1h"`
	BotServiceToken       string        `env:"BOT_SERVICE_TOKEN,required"`
	IngestionServiceToken string        `env:"INGESTION_SERVICE_TOKEN,required"`
	CORSAllowedOrigin     string        `env:"CORS_ALLOWED_ORIGIN,required"`
}
