package businesstype

import (
	"context"
	"errors"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/pgerrors"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/business_type/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/business_type/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/business_type/query"
)

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, businessType *domain.BusinessType) (*domain.BusinessType, error) {
	businessTypeModel := dto.ToModel(*businessType)
	err := r.db.QueryRow(
		ctx,
		query.Upsert,
		businessTypeModel.InfraTypeID,
	).Scan(&businessTypeModel.ID)
	if err != nil {
		return nil, translateWriteError(err)
	}

	return r.GetByID(ctx, businessTypeModel.ID)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BusinessType, error) {
	businessTypeModel := model.BusinessType{}
	err := scanBusinessType(r.db.QueryRow(ctx, query.GetByID, id), &businessTypeModel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrBusinessTypeNotFound
		}
		return nil, err
	}

	return dto.ToDomain(businessTypeModel), nil
}

func (r *repository) GetAll(ctx context.Context) ([]domain.BusinessType, error) {
	rows, err := r.db.Query(ctx, query.GetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	businessTypes := make([]domain.BusinessType, 0)
	for rows.Next() {
		businessTypeModel := model.BusinessType{}
		if err := scanBusinessType(rows, &businessTypeModel); err != nil {
			return nil, err
		}
		businessTypes = append(businessTypes, *dto.ToDomain(businessTypeModel))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return businessTypes, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	var deletedID uuid.UUID
	err := r.db.QueryRow(ctx, query.Delete, id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrBusinessTypeNotFound
		}
		return err
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanBusinessType(row scanner, businessType *model.BusinessType) error {
	return row.Scan(
		&businessType.ID,
		&businessType.InfraTypeID,
		&businessType.InfraType.ID,
		&businessType.InfraType.Slug,
		&businessType.InfraType.Name,
		&businessType.InfraType.Weight,
		&businessType.InfraType.MaxRadius,
		&businessType.CreatedAt,
		&businessType.UpdatedAt,
	)
}

func translateWriteError(err error) error {
	if pgerrors.IsForeignKeyViolation(err) {
		return errs.ErrInfraTypeNotFound
	}
	return err
}
