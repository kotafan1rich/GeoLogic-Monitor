package errs

import "errors"

var (
	ErrTrackedLocationNotFound      = errors.New("tracked location not found")
	ErrTrackedLocationAlreadyExists = errors.New("tracked location already exists")
)
