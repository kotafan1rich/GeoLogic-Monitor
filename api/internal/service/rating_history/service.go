package ratinghistory

import (
	"context"
	"log/slog"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type RatingRepository interface {
	GetHistory(ctx context.Context, trackedLocationID uuid.UUID, months uint) (*domain.LocationRatingHistory, error)
}

type TrackedLocationService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TrackedLocation, error)
}

type service struct {
	log                    *logger.Logger
	ratingRepo             RatingRepository
	trackedLocationService TrackedLocationService
}

func NewService(
	log *logger.Logger,
	ratingRepo RatingRepository,
	trackedLocationService TrackedLocationService,
) *service {
	return &service{
		log:                    log,
		ratingRepo:             ratingRepo,
		trackedLocationService: trackedLocationService,
	}
}

func (s *service) GetHistory(
	ctx context.Context,
	maxUserID int64,
	trackedLocationID uuid.UUID,
	months uint,
) (*domain.LocationRatingHistory, error) {
	trackedLocation, err := s.trackedLocationService.GetByID(ctx, trackedLocationID)
	if err != nil {
		return nil, err
	}

	if trackedLocation.User.MaxUserID != maxUserID {
		return nil, app.Wrap(errs.ErrTrackedLocationNotFound, app.ErrNotFound)
	}

	history, err := s.ratingRepo.GetHistory(ctx, trackedLocationID, months)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get location rating history",
			slog.String("tracked_location_id", trackedLocationID.String()),
			slog.Uint64("months", uint64(months)),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return history, nil
}
