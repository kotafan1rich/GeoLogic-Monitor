package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v7"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool       *pgxpool.Pool
	log        *slog.Logger
	newBackOff func() backoff.BackOff
	maxRetries uint
}

func New(
	ctx context.Context,
	dsn string,
	log *slog.Logger,
	minConns int,
	maxConns int,
	maxConnIdleTime time.Duration,
	maxConnLifetime time.Duration,
) (*DB, error) {
	if log == nil {
		return nil, ErrInvalidLogger
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseDatabaseConfig, err)
	}

	cfg.MinConns = int32(minConns)
	cfg.MaxConns = int32(maxConns)
	cfg.MaxConnIdleTime = maxConnIdleTime
	cfg.MaxConnLifetime = maxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInitPgxPool, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPingDatabase, err)
	}

	return &DB{
		pool:       pool,
		log:        log,
		newBackOff: func() backoff.BackOff { return backoff.NewExponentialBackOff() },
	}, nil
}

func MustNew(
	ctx context.Context,
	dsn string,
	log *slog.Logger,
	minConns int,
	maxConns int,
	maxConnIdleTime time.Duration,
	maxConnLifetime time.Duration,
) *DB {
	const op = "postgresql.MustNew"

	p, err := New(ctx, dsn, log, minConns, maxConns, maxConnIdleTime, maxConnLifetime)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init db: %v", op, err))
	}

	return p
}

func (db *DB) Ping(ctx context.Context) error { return db.pool.Ping(ctx) }

func (db *DB) Close() error {
	db.pool.Close()
	return nil
}
