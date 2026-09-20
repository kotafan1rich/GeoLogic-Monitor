package rating

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type RatingRepository interface {
	Create(ctx context.Context, locationRating *domain.LocationRating) (*domain.LocationRating, error)
}

type Calculator interface {
	Calculate(context.Context, domain.LocationFeatures) (*domain.CalculatedRating, error)
}

type Service interface {
	Calculate(ctx context.Context, features domain.LocationFeatures) (*domain.CalculatedRating, error)
	Create(
		ctx context.Context,
		trackedLocationID uuid.UUID,
		rating *domain.CalculatedRating,
	) (*domain.LocationRating, error)
}

type service struct {
	log        *logger.Logger
	calculator Calculator
	ratingRepo RatingRepository
}

func NewService(
	log *logger.Logger,
	calculator Calculator,
	ratingRepo RatingRepository,
) Service {
	return &service{
		log:        log,
		calculator: calculator,
		ratingRepo: ratingRepo,
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

func (s *service) Create(
	ctx context.Context,
	trackedLocationID uuid.UUID,
	rating *domain.CalculatedRating,
) (*domain.LocationRating, error) {
	locationRating, err := domain.NewLocationRating(trackedLocationID, *rating)
	if err != nil {
		return nil, err
	}
	locationRating, err = s.ratingRepo.Create(ctx, locationRating)
	if err != nil {
		if errors.Is(err, errs.ErrTrackedLocationNotFound) {
			return nil, app.Wrap(err, app.ErrNotFound)
		}
		s.log.ErrorContext(ctx,
			"failed to create location rating",
			slog.String("tracked_location_id", trackedLocationID.String()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return locationRating, nil
}
