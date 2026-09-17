package errs

import "errors"

var (
	ErrInvalidMaxUserID = errors.New("invalid MAX user ID")
	ErrUserNotFound     = errors.New("user not found")
)
