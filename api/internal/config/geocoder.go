package config

import "time"

type Geocoder struct {
	BaseURL string        `env:"GEOCODER_BASE_URL" envDefault:"https://geocode.gate.petersburg.ru"`
	Timeout time.Duration `env:"GEOCODER_TIMEOUT" envDefault:"5s"`
}
