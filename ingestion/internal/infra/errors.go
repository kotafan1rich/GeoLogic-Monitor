package infra

import "errors"

var (
	ErrInvalidAuthToken      = errors.New("auth token is empty")
	ErrInvalidRequestTimeout = errors.New("request timeout must be positive")
	ErrInvalidHTTPClient     = errors.New("HTTP client is nil")
	ErrRateLimit             = errors.New("rate limiter rejected request")
	ErrMarshalData           = errors.New("failed to marshal data")
	ErrUnmarshalData         = errors.New("failed to unmarshal data")
	ErrBuildRequest          = errors.New("failed to build request")
	ErrDoRequest             = errors.New("failed to do request")
	ErrReadResponse          = errors.New("failed to read response")
	ErrUnexpectedStatus      = errors.New("unexpected response status")
	ErrBodyTooLarge          = errors.New("response body exceeds limit")
	ErrEmptyResponse         = errors.New("response body is empty")
)
