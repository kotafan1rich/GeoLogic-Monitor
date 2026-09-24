package errs

import "errors"

var (
	ErrInvalidEventProvider   = errors.New("invalid event provider")
	ErrInvalidEventExternalID = errors.New("invalid event external ID")
	ErrEventNotFound          = errors.New("event not found")
)
