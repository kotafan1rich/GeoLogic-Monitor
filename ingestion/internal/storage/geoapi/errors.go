package geoapi

import "errors"

var (
	ErrInvalidRequester   = errors.New("requester is nil")
	ErrParseURL           = errors.New("failed to parse URL")
	ErrInvalidURL         = errors.New("URL is nil or malformed")
	ErrInvalidClient      = errors.New("geo-api client is nil")
	ErrPutInfraType       = errors.New("failed to put infra type")
	ErrEmptyInfraTypes    = errors.New("infra types list is empty")
	ErrInvalidInfraType   = errors.New("infra type slug is empty")
	ErrDuplicateInfraType = errors.New("duplicate infra type slug")
	ErrUnknownInfraType   = errors.New("unknown infra type slug")
	ErrEmptyTypeID        = errors.New("infra type id is empty in response")
	ErrEmptyAddressList   = errors.New("address list is empty")
	ErrEmptyAddress       = errors.New("address is empty in response")
)
