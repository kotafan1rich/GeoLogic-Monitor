package config

import (
	"log/slog"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	HttpServer HttpServer
	Logging    Logging
	Security   Security
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	cfg := &Config{}

	err = env.Parse(cfg)
	if err != nil {
		slog.Error("error load config", "err", err)
		return nil, err
	}

	return cfg, nil
}
