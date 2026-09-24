package job

import (
	"context"
	"log/slog"
	"time"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/logger"
)

const (
	InfraName  = "infra"
	EventsName = "events"
)

type parseFunc func(ctx context.Context) (int, error)

type datasetsFunc func() []dataset

type Job struct {
	name     string
	schedule string
	log      *slog.Logger
	sources  []datasetsFunc
}

func New(name, schedule string, log *slog.Logger, sources ...datasetsFunc) *Job {
	return &Job{name: name, schedule: schedule, log: log, sources: sources}
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) Schedule() string {
	return j.schedule
}

func (j *Job) Run(ctx context.Context) {
	datasets := make([]dataset, 0, len(j.sources))
	for _, source := range j.sources {
		datasets = append(datasets, source()...)
	}

	runDatasets(ctx, j.log, j.name, datasets)
}

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

	var succeeded, failed int

	for i, ds := range datasets {
		if err := ctx.Err(); err != nil {
			log.WarnContext(ctx, "parse run interrupted",
				slog.Int("processed", i),
				slog.Int("datasets", len(datasets)),
				slog.Any("error", err),
			)

			break
		}

		dsLog := log.With(
			slog.String("dataset", ds.name),
			slog.String("source", ds.source),
		)

		dsLog.DebugContext(ctx, "dataset parsing started",
			slog.String("kind", ds.kind.String()),
			slog.String("target", ds.target),
		)

		dsStart := time.Now()

		count, err := ds.parse(ctx)
		if err != nil {
			failed++

			dsLog.ErrorContext(ctx, "dataset parsing failed",
				slog.String("kind", ds.kind.String()),
				slog.String("target", ds.target),
				slog.Duration("duration", time.Since(dsStart)),
				slog.Any("error", err),
			)

			continue
		}

		succeeded++

		dsLog.InfoContext(ctx, "dataset parsed",
			slog.Int("items", count),
			slog.Duration("duration", time.Since(dsStart)),
		)
	}

	attrs := []any{
		slog.Int("succeeded", succeeded),
		slog.Int("failed", failed),
		slog.Duration("duration", time.Since(runStart)),
	}

	if failed > 0 {
		log.ErrorContext(ctx, "parse run finished with errors", attrs...)

		return
	}

	log.InfoContext(ctx, "parse run finished", attrs...)
}
