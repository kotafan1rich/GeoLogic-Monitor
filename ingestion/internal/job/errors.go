package job

import "errors"

var (
	ErrResolveType = errors.New("failed to resolve infra type id")
	ErrStoreFailed = errors.New("failed to store data")

	ErrConnectBusiness = errors.New("failed to connect business type with infra")
)
