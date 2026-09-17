package config

import "errors"

var (
	ErrReadEnv    = errors.New("failed to read environment variables")
	ErrReadConfig = errors.New("failed to read config file")
)
