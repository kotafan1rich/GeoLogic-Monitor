package geoapi

import "errors"

var (
	ErrInvalidHTTPClient  = errors.New("HTTP client is nil")
	ErrParseURL           = errors.New("failed to parse URL")
	ErrInvalidURL         = errors.New("URL is nil or malformed")
	ErrMarshalData        = errors.New("failed to marshal data")
	ErrBuildRequest       = errors.New("failed to build request")
	ErrCreateRequest      = errors.New("failed to create request")
	ErrReadResponse       = errors.New("failed to read response")
	ErrUnexpectedStatus   = errors.New("unexpected response status")
	ErrUnmarshalData      = errors.New("failed to unmarshal data")
	ErrInvalidClient      = errors.New("geo-api client is nil")
	ErrPutInfraType       = errors.New("failed to put infra type")
	ErrEmptyInfraTypes    = errors.New("infra types list is empty")
	ErrInvalidInfraType   = errors.New("infra type slug is empty")
	ErrDuplicateInfraType = errors.New("duplicate infra type slug")
	ErrUnknownInfraType   = errors.New("unknown infra type slug")
	ErrEmptyTypeID        = errors.New("infra type id is empty in response")
)
