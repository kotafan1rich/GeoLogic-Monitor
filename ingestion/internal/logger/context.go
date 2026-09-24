package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
)

type ctxKey struct{}

const (
	idLen = 8
	idKey = "run_id"
)

func NewRunID() string {
	b := make([]byte, idLen)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func WithRunID(ctx context.Context, runID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, runID)
}

func RunID(ctx context.Context) string {
	if id, ok := ctx.Value(ctxKey{}).(string); ok {
		return id
	}
	return ""
}

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := RunID(ctx); id != "" {
		r.AddAttrs(slog.String(idKey, id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{Handler: h.Handler.WithGroup(name)}
}
