package infraobject

import (
	"context"
	"errors"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/pgerrors"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_object/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_object/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_object/query"
	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, infraObject *domain.InfraObject) (*domain.InfraObject, error) {
	infraObjectModel := dto.ToModel(*infraObject)
	var id uuid.UUID
	err := r.db.QueryRow(
		ctx,
		query.Upsert,
		infraObjectModel.ExternalID,
		infraObjectModel.TypeID,
		infraObjectModel.Location,
		infraObjectModel.Address,
		infraObjectModel.Name,
	).Scan(&id)
	if err != nil {
		if pgerrors.IsForeignKeyViolation(err) {
			return nil, errs.ErrInfraTypeNotFound
		}
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraObject, error) {
	infraObjectModel := model.InfraObject{}
	err := scanInfraObject(r.db.QueryRow(ctx, query.GetByID, id), &infraObjectModel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrInfraObjectNotFound
		}
		return nil, err
	}

	return dto.ToDomain(infraObjectModel), nil
}

func (r *repository) Near(ctx context.Context, geopoint *domain.GeoPoint) ([]*domain.InfraObject, error) {
	location := basemodel.GeoPoint(*geopoint)
	rows, err := r.db.Query(ctx, query.Near, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infraObjects := make([]*domain.InfraObject, 0)
	for rows.Next() {
		infraObjectModel := model.InfraObject{}
		if err := scanInfraObject(rows, &infraObjectModel); err != nil {
			return nil, err
		}
		infraObjects = append(infraObjects, dto.ToDomain(infraObjectModel))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return infraObjects, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanInfraObject(row scanner, infraObject *model.InfraObject) error {
	return row.Scan(
		&infraObject.ID,
		&infraObject.ExternalID,
		&infraObject.TypeID,
		&infraObject.Location,
		&infraObject.Address,
		&infraObject.Name,
		&infraObject.Type.ID,
		&infraObject.Type.Slug,
		&infraObject.Type.Name,
		&infraObject.Type.Weight,
		&infraObject.Type.MaxRadius,
		&infraObject.CreatedAt,
		&infraObject.UpdatedAt,
	)
}
