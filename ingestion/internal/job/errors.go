package job

import "errors"

var (
	ErrInvalidURL  = errors.New("URL is nil or malformed")
	ErrResolveType = errors.New("failed to resolve infra type id")
	ErrStoreFailed = errors.New("failed to store data")
)
