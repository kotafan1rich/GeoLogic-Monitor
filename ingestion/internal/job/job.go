package job

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/logger"
)

type parseFunc func(ctx context.Context, baseURL *url.URL) (int, error)

type dataset struct {
	name    string
	source  string
	baseURL *url.URL
	parse   parseFunc
}

func collect[T any](fn func(context.Context, *url.URL) ([]T, error)) parseFunc {
	return func(ctx context.Context, baseURL *url.URL) (int, error) {
		data, err := fn(ctx, baseURL)
		if err != nil {
			return 0, err
		}

		_ = data

		return len(data), nil
	}
}

func runDatasets(ctx context.Context, log *slog.Logger, jobName string, datasets []dataset) {
	ctx = logger.WithRunID(ctx, logger.NewRunID())
	log = log.With(slog.String("job", jobName))

	runStart := time.Now()

	log.InfoContext(ctx, "parse run started", slog.Int("datasets", len(datasets)))

	var (
		g         errgroup.Group
		succeeded atomic.Int64
		failed    atomic.Int64
	)

	for _, ds := range datasets {
		ds := ds

		g.Go(func() error {
			if ds.baseURL == nil {
				failed.Add(1)

				log.ErrorContext(ctx, "base url is not configured for source",
					slog.String("dataset", ds.name),
					slog.String("source", ds.source),
				)

				return fmt.Errorf("%s [%s]: %w", ds.name, ds.source, ErrInvalidURL)
			}

			log.DebugContext(ctx, "dataset parsing started",
				slog.String("dataset", ds.name),
				slog.String("source", ds.source),
				slog.String("base_url", ds.baseURL.String()),
			)

			dsStart := time.Now()

			count, err := ds.parse(ctx, ds.baseURL)
			if err != nil {
				failed.Add(1)

				log.ErrorContext(ctx, "dataset parsing failed",
					slog.String("dataset", ds.name),
					slog.String("source", ds.source),
					slog.Duration("duration", time.Since(dsStart)),
					slog.Any("error", err),
				)

				return fmt.Errorf("%s: %w", ds.name, err)
			}

			succeeded.Add(1)

			log.InfoContext(ctx, "dataset parsed",
				slog.String("dataset", ds.name),
				slog.String("source", ds.source),
				slog.Int("items", count),
				slog.Duration("duration", time.Since(dsStart)),
			)

			return nil
		})
	}

	err := g.Wait()

	attrs := []any{
		slog.Int64("succeeded", succeeded.Load()),
		slog.Int64("failed", failed.Load()),
		slog.Duration("duration", time.Since(runStart)),
	}

	if err != nil {
		log.ErrorContext(ctx, "parse run finished with errors", append(attrs, slog.Any("error", err))...)
		return
	}

	log.InfoContext(ctx, "parse run finished", attrs...)
}
