package job

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/logger"
)

type parseFunc func(ctx context.Context) (int, error)

type dataset struct {
	name   string
	source string
	kind   a.Kind
	target string
	parse  parseFunc
}

func ds[S a.Source, T any](
	name, source string, src S, fetch func(context.Context, S) ([]T, error), store storeFunc[T],
) dataset {
	return dataset{
		name:   name,
		source: source,
		kind:   src.Kind(),
		target: src.String(),
		parse:  collect(src, fetch, store),
	}
}

func collect[S a.Source, T any](
	src S, fetch func(context.Context, S) ([]T, error), store storeFunc[T],
) parseFunc {
	return func(ctx context.Context) (int, error) {
		if err := src.Validate(); err != nil {
			return 0, err
		}

		data, err := fetch(ctx, src)
		if err != nil {
			return 0, err
		}

		if err := store(ctx, data); err != nil {
			return len(data), err
		}

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
		errs      []error
	)

	for _, ds := range datasets {
		ds := ds

		g.Go(func() error {
			log.DebugContext(ctx, "dataset parsing started",
				slog.String("dataset", ds.name),
				slog.String("source", ds.source),
				slog.String("kind", ds.kind.String()),
				slog.String("target", ds.target),
			)

			dsStart := time.Now()

			count, err := ds.parse(ctx)
			if err != nil {
				failed.Add(1)
				errs = append(errs, err)

				log.ErrorContext(ctx, "dataset parsing failed",
					slog.String("dataset", ds.name),
					slog.String("source", ds.source),
					slog.String("kind", ds.kind.String()),
					slog.String("target", ds.target),
					slog.Duration("duration", time.Since(dsStart)),
					slog.Any("error", err),
				)

				return nil
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

	_ = g.Wait()

	attrs := []any{
		slog.Int64("succeeded", succeeded.Load()),
		slog.Int64("failed", failed.Load()),
		slog.Duration("duration", time.Since(runStart)),
	}

	if len(errs) > 0 {
		log.ErrorContext(
			ctx,
			"parse run finished with errors",
			append(attrs, slog.Any("error", errors.Join(errs...)))...,
		)
		return
	}

	log.InfoContext(ctx, "parse run finished", attrs...)
}
