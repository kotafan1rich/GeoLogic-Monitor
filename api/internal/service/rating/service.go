package rating

import (
	"context"
	"log/slog"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type Calculator interface {
	Calculate(context.Context, domain.LocationFeatures) (*domain.CalculatedRating, error)
}

type service struct {
	calculator Calculator
	log        *logger.Logger
}

func NewService(log *logger.Logger, calculator Calculator) *service {
	return &service{
		calculator: calculator,
		log:        log,
	}
}

func (s *service) Calculate(
	ctx context.Context,
	features domain.LocationFeatures,
) (*domain.CalculatedRating, error) {
	calculatedRating, err := s.calculator.Calculate(ctx, features)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to calculate location rating",
			slog.String("business_type_id", features.BusinessTypeID.String()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return calculatedRating, nil
}
