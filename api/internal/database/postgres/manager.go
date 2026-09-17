package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
)

type manager struct {
	pool   database.DBTX
	logger txLogger
}

type txLogger interface {
	ErrorContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
}

type errorWrapper interface {
	WrapError(ctx context.Context, err error) error
}

func NewManager(pool database.DBTX, logger txLogger) *manager {
	return &manager{
		pool:   pool,
		logger: logger,
	}
}

func (m *manager) wrapError(ctx context.Context, err error) error {
	wrapper, ok := m.logger.(errorWrapper)
	if !ok {
		return err
	}
	return wrapper.WrapError(ctx, err)
}

func (m *manager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return m.wrapError(ctx, fmt.Errorf("begin tx: %w", err))
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			m.logger.ErrorContext(ctx,
				"failed to rollback tx",
				"error", err,
			)
		}
		m.logger.InfoContext(ctx, "transaction rollback")
	}()

	if err := fn(withTx(ctx, tx)); err != nil {
		return m.wrapError(ctx, fmt.Errorf("do fn: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return m.wrapError(ctx, fmt.Errorf("commit tx: %w", err))
	}

	return nil
}
