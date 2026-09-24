package errs

import "errors"

var (
	ErrInvalidLat  = errors.New("invalid latitude")
	ErrInvalidLong = errors.New("invalid longitude")
)
