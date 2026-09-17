package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

type Options struct {
	Level  string
	Format string
	Out    io.Writer
}

const (
	LevelDebug = "debug"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelInfo  = "info"

	TextFormat = "text"
	JSONFormat = "json"

	idLen = 8
)

type ctxKey string

const (
	keyLogger ctxKey = "logger"
	keyRunID  ctxKey = "run_id"
)

func New(opts Options) (*Logger, error) {
	var l slog.Level
	switch opts.Level {
	case LevelDebug:
		l = slog.LevelDebug
	case LevelWarn:
		l = slog.LevelWarn
	case LevelError:
		l = slog.LevelError
	case LevelInfo:
		l = slog.LevelInfo
	default:
		return nil, errors.New("invalid log level")
	}

	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	handlerOpts := &slog.HandlerOptions{Level: l}

	var h slog.Handler
	switch opts.Format {
	case TextFormat:
		h = slog.NewTextHandler(out, handlerOpts)
	case JSONFormat:
		h = slog.NewJSONHandler(out, handlerOpts)
	default:
		return nil, errors.New("invalid log format")
	}

	return &Logger{Logger: slog.New(h)}, nil
}

func MustNew(opts Options) *Logger {
	const op = "ingestion.logger.MustNew"

	l, err := New(opts)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init logger: %v", op, err))
	}
	return l
}

func NewRunID() string {
	b := make([]byte, idLen)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func WithRunID(ctx context.Context, l *Logger, runID string) context.Context {
	ctx = context.WithValue(ctx, keyRunID, runID)
	return context.WithValue(ctx, keyLogger, &Logger{
		Logger: l.Logger.With(slog.String(string(keyRunID), runID)),
	})
}

func From(ctx context.Context) *Logger {
	l, ok := ctx.Value(keyLogger).(*Logger)
	if ok && l != nil {
		return l
	}
	return &Logger{Logger: slog.Default()}
}

func RunID(ctx context.Context) string {
	id, ok := ctx.Value(keyRunID).(string)
	if ok {
		return id
	}
	return ""
}
