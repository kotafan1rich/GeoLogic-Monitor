package aggregator

import "errors"

var (
	ErrInvalidURL      = errors.New("URL is nil or malformed")
	ErrInvalidFilePath = errors.New("file path is empty")
)
