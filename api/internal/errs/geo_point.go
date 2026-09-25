package errs

import "errors"

var (
	ErrInvalidLat = errors.New("invalid latitude")
	ErrInvalidLon = errors.New("invalid longitude")
)
