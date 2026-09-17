package config

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"


type Logging struct {
	LogLevel     logger.LogLevel  `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat    logger.LogFormat `env:"LOG_FORMAT" envDefault:"json"` // json or text
	LogAddSource bool             `env:"LOG_ADD_SOURCE" envDefault:"true"`
}
