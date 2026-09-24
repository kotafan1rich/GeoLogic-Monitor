package infratype

import (
	"context"
	"errors"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_type/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_type/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_type/query"
)

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, infraType *domain.InfraType) (*domain.InfraType, error) {
	infraTypeModel := dto.ToModel(*infraType)
	err := scanInfraType(
		r.db.QueryRow(
			ctx,
			query.Upsert,
			infraTypeModel.Slug,
			infraTypeModel.Name,
			infraTypeModel.Weight,
			infraTypeModel.MaxRadius,
		),
		infraTypeModel,
	)
	if err != nil {
		return nil, err
	}

	return dto.ToDomain(*infraTypeModel), nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraType, error) {
	infraTypeModel := model.InfraType{}
	err := scanInfraType(r.db.QueryRow(ctx, query.GetByID, id), &infraTypeModel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrInfraTypeNotFound
		}
		return nil, err
	}

	return dto.ToDomain(infraTypeModel), nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanInfraType(row scanner, infraType *model.InfraType) error {
	return row.Scan(
		&infraType.ID,
		&infraType.Slug,
		&infraType.Name,
		&infraType.Weight,
		&infraType.MaxRadius,
		&infraType.CreatedAt,
		&infraType.UpdatedAt,
	)
}
