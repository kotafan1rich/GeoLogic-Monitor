package errs

import "errors"

var (
	ErrInvalidSlug         = errors.New("invalid slug")
	ErrInvalidName         = errors.New("invalid name")
	ErrInvalidWeight       = errors.New("invalid weight")
	ErrInvalidRadius       = errors.New("invalid radius")
	ErrInvalidAddress      = errors.New("invalid address")
	ErrInvalidExternalID   = errors.New("invalid external ID")
	ErrInfraObjectNotFound = errors.New("infra object not found")
)
