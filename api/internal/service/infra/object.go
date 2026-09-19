package infra

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	domainerrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type InfraRepository interface {
	Upsert(ctx context.Context, infraObject *domain.InfraObject) (*domain.InfraObject, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraObject, error)
	Near(ctx context.Context, geoPoint *domain.GeoPoint) ([]*domain.InfraObject, error)
}

type infraService struct {
	repo InfraRepository
	log  *logger.Logger
}

func NewInfraService(log *logger.Logger, repo InfraRepository) *infraService {
	return &infraService{log: log, repo: repo}
}

func (s *infraService) Upsert(
	ctx context.Context,
	id uuid.UUID,
	typeID uuid.UUID,
	lat float64,
	lng float64,
	address string,
	name *string,
) (*domain.InfraObject, error) {
	geoPoint, err := domain.NewGeoPoint(lat, lng)
	if err != nil {
		return nil, apperrs.ValidationError(err)
	}

	infraObject, err := domain.NewInfraObject(id, typeID, geoPoint, address, name)
	if err != nil {
		return nil, apperrs.ValidationError(err)
	}

	infraObject, err = s.repo.Upsert(ctx, infraObject)
	if err != nil {
		if errors.Is(err, domainerrs.ErrInfraTypeNotFound) {
			return nil, apperrs.Wrap(err, apperrs.ErrNotFound)
		}

		s.log.ErrorContext(ctx,
			"failed to upsert infra object",
			slog.String("id", id.String()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return infraObject, nil
}

func (s *infraService) GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraObject, error) {
	infraObject, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainerrs.ErrInfraObjectNotFound) {
			return nil, apperrs.Wrap(err, apperrs.ErrNotFound)
		}

		s.log.ErrorContext(ctx,
			"failed to get infra object",
			slog.String("id", id.String()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return infraObject, nil
}

func (s *infraService) Near(
	ctx context.Context,
	geoPoint *domain.GeoPoint,
) ([]*domain.InfraObject, error) {
	infraObjects, err := s.repo.Near(ctx, geoPoint)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get nearby infra objects",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return infraObjects, nil
}
