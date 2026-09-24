package osrm

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

func TestWalkingDistancesMapsRepositoryError(t *testing.T) {
	t.Parallel()

	repositoryErr := errors.New("osrm unavailable")
	service := NewService(testLogger(), &fakeRepository{err: repositoryErr})

	_, err := service.WalkingDistances(
		context.Background(),
		&domain.GeoPoint{},
		[]*domain.GeoPoint{{}},
	)
	var appErr *apperrs.Error
	if !errors.As(err, &appErr) || appErr.Code != apperrs.ErrProviderUnavailable.Code {
		t.Fatalf("error: got %v, want provider_unavailable", err)
	}
	if !errors.Is(err, repositoryErr) {
		t.Fatal("repository error is not preserved")
	}
}

type fakeRepository struct {
	distances []*float64
	err       error
	dst       []*domain.GeoPoint
}

func (r *fakeRepository) WalkingDistances(
	_ context.Context,
	_ *domain.GeoPoint,
	dst []*domain.GeoPoint,
) ([]*float64, error) {
	r.dst = dst
	return r.distances, r.err
}

func testLogger() *logger.Logger {
	return logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
}
