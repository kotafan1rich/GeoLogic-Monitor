package job

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const (
	defaultWriteConcurrency = 8

	externalIDSep = ":"
	datasetPart   = 1

	otherDataset = "other"

	kindInfra    = "infra"
	kindEvents   = "events"
	kindBusiness = "business_types"
)

type storeFunc[T any] = func(ctx context.Context, data []T) error

type Writer interface {
	PutInfraObject(ctx context.Context, obj geoapi.InfraObjectInput) (geoapi.InfraObject, error)
	PutEvent(ctx context.Context, obj geoapi.EventInput) (geoapi.Event, error)
	ConnectWithInfra(ctx context.Context, ID string) (geoapi.BusinessType, error)
}

type TypeResolver interface {
	TypeID(ctx context.Context, slug string) (string, error)
}

type store struct {
	log         *slog.Logger
	writer      Writer
	types       TypeResolver
	concurrency int
}

func (s *store) limit() int {
	if s.concurrency <= 0 {
		return defaultWriteConcurrency
	}

	return s.concurrency
}

func (s *store) infra(slug string) storeFunc[geoapi.InfraObjectInput] {
	return func(ctx context.Context, data []geoapi.InfraObjectInput) error {
		if len(data) == 0 {
			return nil
		}

		byObject := perObjectType(slug)

		var (
			typeID string
			err    error
		)

		if !byObject {
			typeID, err = s.types.TypeID(ctx, slug)
			if err != nil {
				return fmt.Errorf("%w [%s]: %v", ErrResolveType, slug, err)
			}
		}

		var (
			g         errgroup.Group
			collector = errs.NewCollector(errs.DefaultSampleSize)
			stored    atomic.Int64
			skipped   int
			start     = time.Now()
		)

		g.SetLimit(s.limit())

		for _, obj := range data {
			if obj.Lat == 0 && obj.Lon == 0 {
				skipped++
				continue
			}

			if byObject {
				dataset := extractDataset(obj.ExternalID)

				typeID, err = s.types.TypeID(ctx, dataset)
				if err != nil {
					return fmt.Errorf("%w [%s]: %v", ErrResolveType, dataset, err)
				}
			}

			obj.TypeID = typeID

			g.Go(func() error {
				if _, err := s.writer.PutInfraObject(ctx, obj); err != nil {
					collector.Add(err)

					return nil
				}

				stored.Add(1)

				return nil
			})
		}

		_ = g.Wait()

		s.logResult(ctx, kindInfra, slug, len(data), int(stored.Load()), skipped, collector, start)

		if n := collector.Total(); n > 0 {
			return fmt.Errorf("%w: %d/%d infra objects [%s]", ErrStoreFailed, n, len(data), slug)
		}

		return nil
	}
}

func (s *store) events(dataset string) storeFunc[geoapi.EventInput] {
	return func(ctx context.Context, data []geoapi.EventInput) error {
		if len(data) == 0 {
			return nil
		}

		var (
			g         errgroup.Group
			collector = errs.NewCollector(errs.DefaultSampleSize)
			stored    atomic.Int64
			skipped   int
			start     = time.Now()
		)

		g.SetLimit(s.limit())

		for _, obj := range data {
			if obj.Lat == 0 && obj.Lon == 0 || obj.Date.IsZero() {
				skipped++
				continue
			}

			g.Go(func() error {
				if _, err := s.writer.PutEvent(ctx, obj); err != nil {
					collector.Add(err)

					return nil
				}

				stored.Add(1)

				return nil
			})
		}

		_ = g.Wait()

		s.logResult(ctx, kindEvents, dataset, len(data), int(stored.Load()), skipped, collector, start)

		if n := collector.Total(); n > 0 {
			return fmt.Errorf("%w: %d/%d events [%s]", ErrStoreFailed, n, len(data), dataset)
		}

		return nil
	}
}

func (s *store) connectBusiness(ctx context.Context, slugs []string) error {
	if len(slugs) == 0 {
		return nil
	}

	var (
		g         errgroup.Group
		collector = errs.NewCollector(errs.DefaultSampleSize)
		connected atomic.Int64
		start     = time.Now()
	)

	g.SetLimit(s.limit())

	for _, slug := range slugs {
		g.Go(func() error {
			typeID, err := s.types.TypeID(ctx, slug)
			if err != nil {
				collector.Add(fmt.Errorf("%w [%s]: %v", ErrResolveType, slug, err))

				return nil
			}

			if _, err := s.writer.ConnectWithInfra(ctx, typeID); err != nil {
				collector.Add(fmt.Errorf("%w [%s]: %v", ErrConnectBusiness, slug, err))

				return nil
			}

			connected.Add(1)

			return nil
		})
	}

	_ = g.Wait()

	s.logResult(ctx, kindBusiness, datasetBusiness, len(slugs), int(connected.Load()), 0, collector, start)

	if n := collector.Total(); n > 0 {
		return fmt.Errorf("%w: %d/%d business types", ErrConnectBusiness, n, len(slugs))
	}

	return nil
}

func (s *store) logResult(
	ctx context.Context,
	kind, dataset string,
	total, stored, skipped int,
	collector *errs.Collector,
	start time.Time,
) {
	attrs := []any{
		slog.String("kind", kind),
		slog.String("dataset", dataset),
		slog.Int("total", total),
		slog.Int("stored", stored),
		slog.Int("skipped", skipped),
		slog.Int("failed", collector.Total()),
		slog.Int("concurrency", s.limit()),
		slog.Duration("duration", time.Since(start)),
	}

	if collector.Total() > 0 {
		s.log.ErrorContext(ctx, "storing finished with errors", append(attrs, collector.Attr())...)

		return
	}

	if skipped > 0 {
		s.log.WarnContext(ctx, "storing finished, some records skipped", attrs...)

		return
	}

	s.log.DebugContext(ctx, "storing finished", attrs...)
}

func perObjectType(slug string) bool {
	return slug == datasetKidsPlace || slug == datasetBusiness
}

func extractDataset(externalID string) string {
	raw := strings.Split(externalID, externalIDSep)
	if len(raw) <= datasetPart {
		return otherDataset
	}

	return raw[datasetPart]
}
