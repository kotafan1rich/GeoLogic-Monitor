package pgerrors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolationCode     = "23505"
	ForeignKeyViolationCode = "23503"
)

func IsUniqueViolation(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == UniqueViolationCode {
			return true
		}
	}

	return false
}

func IsForeignKeyViolation(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == ForeignKeyViolationCode {
			return true
		}
	}

	return false
}
