package trackedlocation

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	domainerrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/calculate"
)

type InfraService interface {
	Near(ctx context.Context, geoPoint *domain.GeoPoint) ([]*domain.InfraObject, error)
}

type BusinessTypeService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BusinessType, error)
}

type OSRMService interface {
	WalkingDistances(
		ctx context.Context,
		src *domain.GeoPoint,
		dst []*domain.GeoPoint,
	) ([]*float64, error)
}

type TrackedLocationRepository interface {
	Create(ctx context.Context, location *domain.TrackedLocation) (*domain.TrackedLocation, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TrackedLocation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.TrackedLocation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllForMonitoring(ctx context.Context) ([]domain.MonitoringLocation, error)
}

type RatingService interface {
	Calculate(ctx context.Context, features domain.LocationFeatures) (*domain.CalculatedRating, error)
	Create(
		ctx context.Context,
		trackedLocationID uuid.UUID,
		rating *domain.CalculatedRating,
	) (*domain.LocationRating, error)
}

type service struct {
	repo                TrackedLocationRepository
	osrmService         OSRMService
	infraServie         InfraService
	businessTypeService BusinessTypeService
	ratingService       RatingService

	txManager database.TxManager
	log       *logger.Logger
}

func NewTrackedLocationService(
	repo TrackedLocationRepository,
	osrmService OSRMService,
	infraServie InfraService,
	businessTypeService BusinessTypeService,
	ratingService RatingService,
	txManager database.TxManager,
	log *logger.Logger,
) *service {
	return &service{
		repo:                repo,
		osrmService:         osrmService,
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
	name string,
	address string,
	lat float64,
	lng float64,
) (*domain.TrackedLocationRating, error) {
	if strings.TrimSpace(name) == "" {
		return nil, apperrs.ValidationError(domainerrs.ErrInvalidName)
	}
	if strings.TrimSpace(address) == "" {
		return nil, apperrs.ValidationError(domainerrs.ErrInvalidAddress)
	}

	geoPoint, err := domain.NewGeoPoint(lat, lng)
	if err != nil || geoPoint == nil {
		return nil, apperrs.ValidationError(err)
	}

	var calculatedRating *domain.CalculatedRating
	var createdLocation *domain.TrackedLocation

	err = s.txManager.WithTx(ctx, func(ctx context.Context) error {
		location := domain.NewTrackedLocation(userID, businessTypeID, name, address, geoPoint)
		location, err := s.repo.Create(ctx, location)
		if err != nil {
			if errors.Is(err, domainerrs.ErrBusinessTypeNotFound) {
				return apperrs.Wrap(err, apperrs.ErrNotFound)
			}
			s.log.ErrorContext(ctx,
				"failed to create tracked location",
				slog.String("user_id", userID.String()),
				slog.String("error", err.Error()),
			)
			return err
		}

		ratingResult, err := s.calculateRating(ctx, location)
		if err != nil {
			return err
		}
		if _, err := s.ratingService.Create(ctx, location.ID, ratingResult); err != nil {
			return err
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

func (s *service) DeleteForUser(ctx context.Context, id, userID uuid.UUID) error {
	location, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if location.UserID != userID {
		return apperrs.Wrap(domainerrs.ErrTrackedLocationNotFound, apperrs.ErrNotFound)
	}

	return s.Delete(ctx, id)
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

func (s *service) Recalculate(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	location, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	rating, err := s.calculateRating(ctx, location)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err = s.ratingService.Create(ctx, location.ID, rating)
	return err
}

func (s *service) RecalculateAll(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	locations, err := s.GetAllForMonitoring(ctx)
	if err != nil {
		return err
	}
	for _, location := range locations {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := s.Recalculate(ctx, location.ID); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			s.log.ErrorContext(ctx, "failed to recalculate tracked location rating",
				slog.String("tracked_location_id", location.ID.String()),
				slog.String("error", err.Error()),
			)
		}
	}
	return ctx.Err()
}

func (s *service) calculateRating(ctx context.Context, location *domain.TrackedLocation) (*domain.CalculatedRating, error) {
	infraNear, err := s.infraServie.Near(ctx, &location.GeoPoint)
	if err != nil {
		return nil, err
	}

	dst := make([]*domain.GeoPoint, len(infraNear))
	for i := range infraNear {
		dst[i] = &infraNear[i].GeoPoint
	}

	distances, err := s.osrmService.WalkingDistances(ctx, &location.GeoPoint, dst)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get walking distances for location",
			slog.String("user_id", location.UserID.String()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	if len(distances) != len(infraNear) {
		return nil, errors.New("unexpected walking distance count")
	}

	businessType, err := s.businessTypeService.GetByID(ctx, location.BusinessTypeID)
	if err != nil {
		return nil, err
	}

	infraWithDistance := make([]*domain.InfraObjectDistance, 0, len(infraNear))
	for i, distance := range distances {
		if distance == nil || *distance > float64(infraNear[i].Type.MaxRadius) {
			continue
		}
		infraWithDistance = append(infraWithDistance, &domain.InfraObjectDistance{
			Object:         infraNear[i],
			DistanceMeters: *distance,
		})
	}

	features := calculate.BuildLocationFeatures(businessType, infraWithDistance)

	ratingResult, err := s.ratingService.Calculate(ctx, *features)
	if err != nil {
		return nil, apperrs.Wrap(err, apperrs.ErrProviderUnavailable)
	}
	return ratingResult, nil
}
