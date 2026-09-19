package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type Logger struct {
	*slog.Logger
}

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type LogFormat string

const (
	FormatJson LogFormat = "json"
	FormatText LogFormat = "text"
)

func New(level LogLevel, format LogFormat, addSource bool, output io.Writer) *Logger {
	var slogLevel slog.Level
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		panic("bad log level")
	}

	writer := output
	if writer == nil {
		writer = os.Stdout
	}
	opts := &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: addSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.RFC3339))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	switch format {
	case FormatJson:
		handler = slog.NewJSONHandler(writer, opts)
	case FormatText:
		handler = slog.NewTextHandler(writer, opts)
	default:
		panic("bad log format")
	}

	logger := slog.New(handler)

	logger.Info(
		"Created logger",
		slog.String("level", string(level)),
		slog.String("format", string(format)),
		slog.Bool("addSource", addSource),
		slog.Any("output", output),
	)
	return &Logger{
		Logger: logger,
	}
}

func Default() *Logger {
	return New(LevelInfo, FormatText, false, os.Stdout)
}

func (l *Logger) WithLatency(latency time.Duration) *Logger {
	return &Logger{
		Logger: l.With(slog.Duration("latency", latency)),
	}
}

func (l *Logger) WithRequest(status int, method, path, query, ip, userAgent string, latency time.Duration) *Logger {
	return &Logger{
		Logger: l.With(
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", ip),
			slog.String("user-agent", userAgent),
			slog.Duration("latency", latency),
		),
	}
}

func (l *Logger) LogError(ctx context.Context, err error, message string) {
	if err != nil {
		pc := make([]uintptr, 10)
		n := runtime.Callers(2, pc)
		frames := runtime.CallersFrames(pc[:n])

		var file string
		var line int

		if frame, more := frames.Next(); more {
			file = frame.File
			line = frame.Line
		}

		l.ErrorContext(ctx, message,
			slog.String("error", err.Error()),
			slog.String("file", file),
			slog.Int("line", line),
		)
	}
}
