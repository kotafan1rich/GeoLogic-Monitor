package config

import "time"

type Api struct {
	BaseURL string        `env:"API_BASE_URL" envDefault:"http://localhost:8080"`
	Timeout time.Duration `env:"API_TIMEOUT" envDefault:"5s"`
}
