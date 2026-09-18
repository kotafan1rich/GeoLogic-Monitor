package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
)

func TestManagerWithTxCommitsSuccessfulCallback(t *testing.T) {
	t.Parallel()

	tx := &fakeTx{}
	manager := NewManager(&fakeDB{tx: tx}, fakeLogger{})

	err := manager.WithTx(context.Background(), func(ctx context.Context) error {
		contextTx, ok := txFromContext(ctx)
		if !ok {
			t.Fatal("transaction is missing from context")
		}
		if contextTx != tx {
			t.Fatal("unexpected transaction in context")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithTx returned an error: %v", err)
	}
	if tx.commits != 1 {
		t.Fatalf("commit calls: got %d, want 1", tx.commits)
	}
	if tx.rollbacks != 1 {
		t.Fatalf("deferred rollback calls: got %d, want 1", tx.rollbacks)
	}
}

func TestManagerWithTxRollsBackFailedCallback(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("callback failed")
	tx := &fakeTx{}
	manager := NewManager(&fakeDB{tx: tx}, fakeLogger{})

	err := manager.WithTx(context.Background(), func(context.Context) error {
		return expectedErr
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("WithTx error: got %v, want %v", err, expectedErr)
	}
	if tx.commits != 0 {
		t.Fatalf("commit calls: got %d, want 0", tx.commits)
	}
	if tx.rollbacks != 1 {
		t.Fatalf("rollback calls: got %d, want 1", tx.rollbacks)
	}
}

type fakeDB struct {
	tx database.Tx
}

func (db *fakeDB) Exec(context.Context, string, ...any) error {
	return nil
}

func (db *fakeDB) QueryRow(context.Context, string, ...any) database.Row {
	return nil
}

func (db *fakeDB) Query(context.Context, string, ...any) (database.Rows, error) {
	return nil, nil
}

func (db *fakeDB) Begin(context.Context) (database.Tx, error) {
	return db.tx, nil
}

func (db *fakeDB) Close() {}

type fakeTx struct {
	commits   int
	rollbacks int
	closed    bool
}

func (tx *fakeTx) Commit(context.Context) error {
	tx.commits++
	tx.closed = true
	return nil
}

func (tx *fakeTx) Rollback(context.Context) error {
	tx.rollbacks++
	if tx.closed {
		return pgx.ErrTxClosed
	}
	tx.closed = true
	return nil
}

func (tx *fakeTx) Exec(context.Context, string, ...any) error {
	return nil
}

func (tx *fakeTx) QueryRow(context.Context, string, ...any) database.Row {
	return nil
}

func (tx *fakeTx) Query(context.Context, string, ...any) (database.Rows, error) {
	return nil, nil
}

type fakeLogger struct{}

func (fakeLogger) ErrorContext(context.Context, string, ...any) {}

func (fakeLogger) InfoContext(context.Context, string, ...any) {}
