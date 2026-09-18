package digitalspb

import "errors"

var (
	ErrInvalidHTTPClient = errors.New("HTTP client is nil")
	ErrInvalidURLMap     = errors.New("URLs map is empty or malformed")
	ErrBuildRequest      = errors.New("failed to build request")
	ErrCreateRequest     = errors.New("failed to create request")
	ErrReadResponse      = errors.New("failed to read response")
	ErrUnexpectedStatus  = errors.New("unexpected response status")
	ErrUnmarshalData     = errors.New("failed to unmarshal data")
)
