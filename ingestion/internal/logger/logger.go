package logger

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

const (
	LevelDebug = "debug"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelInfo  = "info"
	TextFormat = "text"
	JSONFormat = "json"
)

func New(level, format string, out io.Writer) (*slog.Logger, error) {
	var l slog.Level
	switch level {
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

	if out == nil {
		out = os.Stdout
	}

	handlerOpts := &slog.HandlerOptions{Level: l}

	var h slog.Handler
	switch format {
	case TextFormat:
		h = slog.NewTextHandler(out, handlerOpts)
	case JSONFormat:
		h = slog.NewJSONHandler(out, handlerOpts)
	default:
		return nil, errors.New("invalid log format")
	}

	return slog.New(ContextHandler{Handler: h}), nil
}

func MustNew(level, format string, out io.Writer) *slog.Logger {
	const op = "ingestion.logger.MustNew"

	l, err := New(level, format, out)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init logger: %v", op, err))
	}
	return l
}
