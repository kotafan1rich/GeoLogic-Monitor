package errs

import "errors"

var (
	ErrBusinessTypeNotFound = errors.New("business type not found")
	ErrInfraTypeNotFound    = errors.New("infra type not found")
)
