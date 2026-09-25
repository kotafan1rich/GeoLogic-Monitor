package infratype

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

type InfraTypeRepository interface {
	Upsert(ctx context.Context, infraType *domain.InfraType) (*domain.InfraType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraType, error)
}

type typeService struct {
	repo InfraTypeRepository
	log  *logger.Logger
}

func NewService(log *logger.Logger, repo InfraTypeRepository) *typeService {
	return &typeService{log: log, repo: repo}
}

func (s *typeService) Upsert(
	ctx context.Context,
	slug, name string,
	weight float64,
	maxRadius uint16,
) (*domain.InfraType, error) {
	infraType, err := domain.NewInfraType(slug, name, weight, maxRadius)
	if err != nil {
		return nil, apperrs.ValidationError(err)
	}

	infraType, err = s.repo.Upsert(ctx, infraType)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to upsert infra type",
			slog.String("slug", slug),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return infraType, nil
}

func (s *typeService) GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraType, error) {
	infraType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainerrs.ErrInfraTypeNotFound) {
			return nil, apperrs.Wrap(err, apperrs.ErrNotFound)
		}

		s.log.ErrorContext(ctx,
			"failed to get infra type",
			slog.String("id", id.String()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return infraType, nil
}
