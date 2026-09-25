package maps

import "errors"

var (
	ErrInvalidLogger               = errors.New("logger is nil")
	ErrInvalidRequester            = errors.New("requester is nil")
	ErrInvalidURL                  = errors.New("URL is empty")
	ErrParseURL                    = errors.New("failed to parse URL")
	ErrOpenFile                    = errors.New("failed to open file")
	ErrDecodeData                  = errors.New("failed to decode data")
	ErrInvalidCoordinatesConverter = errors.New("coordinates converter is nil")
)
