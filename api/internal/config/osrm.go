package config

import "time"

type OSRM struct {
	BaseURL string        `env:"OSRM_BASE_URL" envDefault:"http://osrm:5000"`
	Timeout time.Duration `env:"OSRM_TIMEOUT" envDefault:"5s"`
}
