package postgresql

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v7"
	"github.com/jackc/pgx/v5"
)

func (db *DB) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return db.withTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, 1, fn)
}

func (db *DB) WithSerializableTx(ctx context.Context, fn func(context.Context) error) error {
	return db.withTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, db.maxRetries, fn)
}

func (db *DB) withTx(
	ctx context.Context, opts pgx.TxOptions, maxRetries uint, fn func(context.Context) error,
) error {
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	_, err := backoff.Retry(
		ctx,
		func() (struct{}, error) {
			err := db.runTx(ctx, opts, fn)
			if err != nil && !isRetryableErr(err) {
				return struct{}{}, backoff.Permanent(err)
			}
			return struct{}{}, err
		},
		backoff.WithBackOff(db.newBackOff()),
		backoff.WithMaxElapsedTime(0),
		backoff.WithMaxTries(maxRetries),
		backoff.WithNotify(func(err error, d time.Duration) {
			db.log.WarnContext(ctx, "transaction conflict",
				slog.Any("error", err), slog.Duration("next_retry_in", d))
		}),
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) runTx(
	ctx context.Context, opts pgx.TxOptions, fn func(context.Context) error,
) error {
	tx, err := db.pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback(context.WithoutCancel(ctx))
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			db.log.ErrorContext(ctx, "rollback failed", slog.Any("error", err))
		}
	}()

	if err = fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
