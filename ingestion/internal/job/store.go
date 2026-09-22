package job

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"golang.org/x/sync/errgroup"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const infraWriteConcurrency = 8

type storeFunc[T any] = func(ctx context.Context, data []T) error

type InfraWriter interface {
	PutInfraObject(ctx context.Context, obj geoapi.InfraObjectInput) (geoapi.InfraObject, error)
}

type TypeResolver interface {
	TypeID(ctx context.Context, slug string) (string, error)
}

func toInfra(w InfraWriter, types TypeResolver, slug string) storeFunc[geoapi.InfraObjectInput] {
	return func(ctx context.Context, data []geoapi.InfraObjectInput) error {
		if len(data) == 0 {
			return nil
		}

		typeID, err := types.TypeID(ctx, slug)
		if err != nil {
			return fmt.Errorf("%w [%s]: %v", ErrResolveType, slug, err)
		}

		var (
			g      errgroup.Group
			failed atomic.Int64
			errs   []error
		)

		g.SetLimit(infraWriteConcurrency)

		for _, obj := range data {
			obj.TypeID = typeID

			if obj.Lat == 0 && obj.Lon == 0 {
				failed.Add(1)
				continue
			}

			g.Go(func() error {
				if _, err := w.PutInfraObject(ctx, obj); err != nil {
					failed.Add(1)
					errs = append(errs, err)
				}

				return nil
			})
		}

		_ = g.Wait()

		if n := failed.Load(); n > 0 && len(errs) > 0 {
			return fmt.Errorf("%w: %d/%d objects: %v", ErrStoreFailed, n, len(data), errors.Join(errs...))
		}

		return nil
	}
}

func all[T any](stores ...storeFunc[T]) storeFunc[T] {
	return func(ctx context.Context, data []T) error {
		for _, store := range stores {
			if err := store(ctx, data); err != nil {
				return err
			}
		}

		return nil
	}
}

// func discard[T any](_ context.Context, _ []T) error {
// 	return nil
// }
