package osrm

import (
	"context"
	"log/slog"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type OSRMRepository interface {
	WalkingDistances(
		ctx context.Context,
		src *domain.GeoPoint,
		dst []*domain.GeoPoint,
	) ([]*float64, error)
}

type service struct {
	log      *logger.Logger
	osrmRepo OSRMRepository
}

func NewService(log *logger.Logger, osrmRepo OSRMRepository) *service {
	return &service{
		log:      log,
		osrmRepo: osrmRepo,
	}
}

func (s *service) WalkingDistances(
	ctx context.Context,
	src *domain.GeoPoint,
	dst []*domain.GeoPoint,
) ([]*float64, error) {
	distances, err := s.osrmRepo.WalkingDistances(ctx, src, dst)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get walking distance",
			slog.String("error", err.Error()),
		)
		return nil, apperrs.Wrap(err, apperrs.ErrProviderUnavailable)
	}
	return distances, nil
}
