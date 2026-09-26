package postgresql

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidLogger       = errors.New("logger is nil")
	ErrParseDatabaseConfig = errors.New("failed to parse database config")
	ErrInitPgxPool         = errors.New("failed to init pgx pool")
	ErrPingDatabase        = errors.New("failed to ping database")
	ErrBeginTx             = errors.New("failed to begin transaction")
	ErrExecuteTx           = errors.New("failed to execute transaction")
	ErrCommitTx            = errors.New("failed to commit transaction")
)

func isRetryableErr(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == pgerrcode.SerializationFailure ||
		pgErr.Code == pgerrcode.DeadlockDetected
}
