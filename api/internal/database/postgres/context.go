package postgres

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
)

type txKey struct{}

func withTx(ctx context.Context, tx database.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// DBTXFromContext возвращает DBTX, открытую текущим Manager.WithTx, если вызов идёт
// внутри неё; иначе — fallback (обычно *pgxpool.Pool, переданный репозиторию при
// конструировании).
func DBTXFromContext(ctx context.Context, fallback database.Tx) database.Tx {
	if tx, ok := ctx.Value(txKey{}).(database.Tx); ok {
		return tx
	}

	return fallback
}
