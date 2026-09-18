package businesstype

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

type Repository interface {
	Upsert(ctx context.Context, businessType *domain.BusinessType) (*domain.BusinessType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BusinessType, error)
	GetAll(ctx context.Context) ([]domain.BusinessType, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service interface {
	Upsert(ctx context.Context, infraTypeID uuid.UUID) (*domain.BusinessType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BusinessType, error)
	GetAll(ctx context.Context) ([]domain.BusinessType, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
	log  *logger.Logger
}

func NewService(log *logger.Logger, repo Repository) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) Upsert(ctx context.Context, infraTypeID uuid.UUID) (*domain.BusinessType, error) {
	businessType, err := s.repo.Upsert(ctx, domain.NewBusinessType(infraTypeID))
	if err != nil {
		return nil, s.mapError(ctx, err, "failed to upsert business type")
	}

	return businessType, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*domain.BusinessType, error) {
	businessType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, s.mapError(ctx, err, "failed to get business type")
	}

	return businessType, nil
}

func (s *service) GetAll(ctx context.Context) ([]domain.BusinessType, error) {
	businessTypes, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, s.mapError(ctx, err, "failed to get business types")
	}

	return businessTypes, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return s.mapError(ctx, err, "failed to delete business type")
	}

	return nil
}

func (s *service) mapError(ctx context.Context, err error, message string) error {
	if errors.Is(err, domainerrs.ErrBusinessTypeNotFound) ||
		errors.Is(err, domainerrs.ErrInfraTypeNotFound) {
		return apperrs.Wrap(err, apperrs.ErrNotFound)
	}

	s.log.ErrorContext(ctx, message, slog.String("error", err.Error()))
	return err
}
