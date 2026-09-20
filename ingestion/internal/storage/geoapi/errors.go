package geoapi

import "errors"

var (
	ErrInvalidHTTPClient = errors.New("HTTP client is nil")
	ErrParseURL          = errors.New("failed to parse URL")
	ErrInvalidURL        = errors.New("URL is nil or malformed")
	ErrBuildRequest      = errors.New("failed to build request")
	ErrCreateRequest     = errors.New("failed to create request")
	// ErrReadResponse      = errors.New("failed to read response")
	ErrUnexpectedStatus = errors.New("unexpected response status")
)
