package postgres

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
)

type txKey struct{}

func withTx(ctx context.Context, tx database.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func txFromContext(ctx context.Context) (database.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(database.Tx)
	return tx, ok
}
