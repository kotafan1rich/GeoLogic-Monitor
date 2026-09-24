package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func Load(path string) (*Config, error) {
	var cfg Config

	if path == "" {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrReadEnv, err)
		}

		return &cfg, nil
	}

	_ = godotenv.Load()

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadConfig, err)
	}

	return &cfg, nil
}

func MustLoad(path string) {
	const op = "ingestion.config.MustLoad"

	cfg, err := Load(path)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to load config: %v", op, err))
	}
	config = cfg
}
