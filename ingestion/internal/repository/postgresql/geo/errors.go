package geo

import (
	"errors"
)

var (
	ErrInvalidType       = errors.New("invalid type for DBGeoPoint")
	ErrInvalidEWKBLen    = errors.New("invalid EWKB length")
	ErrInsufficientBytes = errors.New("insufficient bytes for coordinates")
)
