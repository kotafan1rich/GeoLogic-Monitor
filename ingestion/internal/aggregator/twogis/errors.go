package twogis

import "errors"

var (
	ErrInvalidLogger            = errors.New("logger is nil")
	ErrInvalidRequester         = errors.New("requester is nil")
	ErrInvalidURL               = errors.New("URL is nil or malformed")
	ErrParseURL                 = errors.New("failed to parse URL")
	ErrInvalidAPIKey            = errors.New("API key is empty")
	ErrUnexpectedProviderStatus = errors.New("unexpected provider status")
	ErrDecodeData               = errors.New("failed to decode data")
)
