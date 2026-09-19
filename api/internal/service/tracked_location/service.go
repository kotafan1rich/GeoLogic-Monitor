package trackedlocation

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	domainerrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/rating"
)

type RatingService interface {
	Calculate(ctx context.Context, features domain.LocationFeatures) (*domain.CalculatedRating, error)
}

type InfraService interface {
	Near(ctx context.Context, geoPoint *domain.GeoPoint) ([]*domain.InfraObject, error)
}

type BusinessTypeService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BusinessType, error)
}

type OSRMRepository interface {
	FilterWalkingDistance(ctx context.Context, src *domain.GeoPoint, dst []*domain.InfraObject) ([]*domain.InfraObjectDistance, error)
}
type TrackedLocationRepository interface {
	Create(ctx context.Context, location *domain.TrackedLocation) (*domain.TrackedLocation, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TrackedLocation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.TrackedLocation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllForMonitoring(ctx context.Context) ([]domain.MonitoringLocation, error)
}

type service struct {
	repo                TrackedLocationRepository
	osrmRepo            OSRMRepository
	infraServie         InfraService
	businessTypeService BusinessTypeService
	ratingService       RatingService

	txManager database.TxManager
	log       *logger.Logger
}

func NewTrackedLocationService(
	repo TrackedLocationRepository,
	osrmRepo OSRMRepository,
	infraServie InfraService,
	businessTypeService BusinessTypeService,
	ratingService RatingService,
	txManager database.TxManager,
	log *logger.Logger,
) *service {
	return &service{
		repo:                repo,
		osrmRepo:            osrmRepo,
		infraServie:         infraServie,
		businessTypeService: businessTypeService,
		ratingService:       ratingService,
		txManager:           txManager,
		log:                 log,
	}
}

func (s *service) Create(
	ctx context.Context,
	userID uuid.UUID,
	businessTypeID uuid.UUID,
	address string,
	lat float64,
	lng float64,
) (*domain.TrackedLocationRating, error) {
	geoPoint, err := domain.NewGeoPoint(lat, lng)
	if err != nil || geoPoint == nil {
		return nil, apperrs.ValidationError(err)
	}

	var calculatedRating *domain.CalculatedRating
	var createdLocation *domain.TrackedLocation

	err = s.txManager.WithTx(ctx, func(ctx context.Context) error {
		location := domain.NewTrackedLocation(userID, businessTypeID, address, geoPoint)
		location, err := s.repo.Create(ctx, location)
		if err != nil {
			s.log.ErrorContext(ctx,
				"failed to create tracked location",
				slog.String("user_id", userID.String()),
				slog.String("error", err.Error()),
			)
			return err
		}

		infraNear, err := s.infraServie.Near(ctx, &location.GeoPoint)
		if err != nil {
			return err
		}

		infraWithDistance, err := s.osrmRepo.FilterWalkingDistance(ctx, &location.GeoPoint, infraNear)
		if err != nil {
			s.log.ErrorContext(ctx,
				"failed to get walking distances for location",
				slog.String("user_id", userID.String()),
				slog.String("error", err.Error()),
			)
			return apperrs.Wrap(err, apperrs.ErrProviderUnavailable)
		}

		businessType, err := s.businessTypeService.GetByID(ctx, location.BusinessTypeID)
		if err != nil {
			return err
		}
		features := rating.BuildLocationFeatures(businessType, infraWithDistance)

		ratingResult, err := s.ratingService.Calculate(ctx, *features)
		if err != nil {
			return apperrs.Wrap(err, apperrs.ErrProviderUnavailable)
		}

		createdLocation = location
		calculatedRating = ratingResult
		return nil
	})

	if err != nil {
		return nil, err
	}

	return domain.NewTrackedLocationRating(createdLocation, calculatedRating), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*domain.TrackedLocation, error) {
	location, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, s.mapError(ctx, err, "failed to get tracked location", id)
	}

	return location, nil
}

func (s *service) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.TrackedLocation, error) {
	locations, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domainerrs.ErrTrackedLocationNotFound) {
			return nil, apperrs.Wrap(err, apperrs.ErrNotFound)
		}

		s.log.ErrorContext(ctx,
			"failed to get tracked locations by user ID",
			slog.String("user_id", userID.String()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return locations, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return s.mapError(ctx, err, "failed to delete tracked location", id)
	}

	return nil
}

func (s *service) GetAllForMonitoring(ctx context.Context) ([]domain.MonitoringLocation, error) {
	locations, err := s.repo.GetAllForMonitoring(ctx)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get tracked locations for monitoring",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return locations, nil
}

func (s *service) mapError(ctx context.Context, err error, message string, id uuid.UUID) error {
	if errors.Is(err, domainerrs.ErrTrackedLocationNotFound) {
		return apperrs.Wrap(err, apperrs.ErrNotFound)
	}

	s.log.ErrorContext(ctx,
		message,
		slog.String("id", id.String()),
		slog.String("error", err.Error()),
	)
	return err
}
