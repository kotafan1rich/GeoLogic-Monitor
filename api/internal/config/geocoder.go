package config

import "time"

type Geocoder struct {
	BaseURL string        `env:"DADATA_BASE_URL" envDefault:"https://suggestions.dadata.ru/suggestions/api/4_1/rs"`
	APIKey  string        `env:"DADATA_API_KEY"`
	Timeout time.Duration `env:"GEOCODER_TIMEOUT" envDefault:"5s"`
}
