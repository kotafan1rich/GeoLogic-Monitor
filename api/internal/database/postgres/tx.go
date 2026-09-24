package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
)

type postgresTx struct {
	tx pgx.Tx
}

func NewTx(tx pgx.Tx) database.Tx {
	return &postgresTx{
		tx: tx,
	}
}

func (t *postgresTx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *postgresTx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *postgresTx) Exec(ctx context.Context, query string, args ...any) error {
	_, err := t.tx.Exec(ctx, query, args...)
	return err
}

func (t *postgresTx) QueryRow(ctx context.Context, query string, args ...any) database.Row {
	return t.tx.QueryRow(ctx, query, args...)
}

func (t *postgresTx) Query(ctx context.Context, query string, args ...any) (database.Rows, error) {
	return t.tx.Query(ctx, query, args...)
}
