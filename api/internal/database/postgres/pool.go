package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
)

type pgxPool struct {
	pool *pgxpool.Pool
}

func NewPool(
	ctx context.Context,
	dsn string,
	minIdleConns int,
	maxConns int,
	maxConnLifetime time.Duration,
) (database.DBTX, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MinIdleConns = int32(minIdleConns)
	poolConfig.MaxConns = int32(maxConns)
	poolConfig.MaxConnLifetime = maxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &pgxPool{pool: pool}, nil
}

func (p *pgxPool) Exec(ctx context.Context, query string, args ...any) error {
	if tx, ok := txFromContext(ctx); ok {
		return tx.Exec(ctx, query, args...)
	}

	_, err := p.pool.Exec(ctx, query, args...)
	return err
}

func (p *pgxPool) QueryRow(ctx context.Context, query string, args ...any) database.Row {
	if tx, ok := txFromContext(ctx); ok {
		return tx.QueryRow(ctx, query, args...)
	}

	return p.pool.QueryRow(ctx, query, args...)
}

func (p *pgxPool) Query(ctx context.Context, query string, args ...any) (database.Rows, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.Query(ctx, query, args...)
	}

	return p.pool.Query(ctx, query, args...)
}

func (p *pgxPool) Begin(ctx context.Context) (database.Tx, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return NewTx(tx), nil
}

func (p *pgxPool) Close() {
	p.pool.Close()
}
