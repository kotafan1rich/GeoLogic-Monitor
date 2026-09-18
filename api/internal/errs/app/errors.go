package app

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	err     error
}

var (
	ErrNotFound = &Error{
		Status:  http.StatusNotFound,
		Code:    "not_found",
		Message: "resource not found",
	}
	ErrUnauthorized = &Error{
		Status:  http.StatusUnauthorized,
		Code:    "unauthorized",
		Message: "unauthorized",
	}
	ErrProviderUnavailable = &Error{
		Status:  http.StatusServiceUnavailable,
		Code:    "provider_unavailable",
		Message: "provider unavailable",
	}
	ErrInternal = &Error{
		Status:  http.StatusInternalServerError,
		Code:    "internal_error",
		Message: "internal server error",
	}
)

func (e *Error) Error() string {
	if e.err == nil {
		return fmt.Sprintf("[%s] %s", e.Code, e.Message)
	}
	return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}

func Wrap(err error, template *Error) *Error {
	return &Error{
		Status:  template.Status,
		Code:    template.Code,
		Message: template.Message,
		err:     err,
	}
}

func ValidationError(err error) *Error {
	return &Error{
		Status:  http.StatusBadRequest,
		Code:    "validation_error",
		Message: err.Error(),
		err:     err,
	}
}
