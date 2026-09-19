package trackedlocations

import (
	"context"
	"errors"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/pgerrors"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/tracked_location/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/tracked_location/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/tracked_location/query"
)

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, location *domain.TrackedLocation) (*domain.TrackedLocation, error) {
	locationModel := dto.ToModel(*location)
	err := r.db.QueryRow(
		ctx,
		query.Create,
		locationModel.UserID,
		locationModel.BusinessTypeID,
		locationModel.Address,
		locationModel.Location,
	).Scan(
		&locationModel.ID,
		&locationModel.UserID,
		&locationModel.BusinessTypeID,
		&locationModel.Address,
		&locationModel.Location,
		&locationModel.CreatedAt,
		&locationModel.UpdatedAt,
	)

	if err != nil {
		if pgerrors.IsUniqueViolation(err) {
			return nil, errs.ErrTrackedLocationAlreadyExists
		}
		if pgerrors.IsForeignKeyViolation(err) {
			return nil, errs.ErrBusinessTypeNotFound
		}
		return nil, err
	}
	return dto.ToDomain(*locationModel), nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.TrackedLocation, error) {
	locationModel := model.TrackedLocation{}
	err := r.db.QueryRow(
		ctx,
		query.GetByID,
		id,
	).Scan(
		&locationModel.ID,
		&locationModel.UserID,
		&locationModel.BusinessTypeID,
		&locationModel.Address,
		&locationModel.Location,
		&locationModel.CreatedAt,
		&locationModel.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrTrackedLocationNotFound
		}
		return nil, err
	}
	return dto.ToDomain(locationModel), nil
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.TrackedLocation, error) {
	rows, err := r.db.Query(
		ctx,
		query.GetByUserID,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]domain.TrackedLocation, 0)
	for rows.Next() {
		locationModel := model.TrackedLocation{}
		if err := rows.Scan(
			&locationModel.ID,
			&locationModel.UserID,
			&locationModel.BusinessTypeID,
			&locationModel.Address,
			&locationModel.Location,
			&locationModel.CreatedAt,
			&locationModel.UpdatedAt,
		); err != nil {
			return nil, err
		}

		locations = append(locations, *dto.ToDomain(locationModel))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	var deletedID uuid.UUID
	err := r.db.QueryRow(
		ctx,
		query.Delete,
		id,
	).Scan(&deletedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrTrackedLocationNotFound
		}
		return err
	}

	return nil
}

func (r *repository) GetAllForMonitoring(ctx context.Context) ([]domain.MonitoringLocation, error) {
	rows, err := r.db.Query(ctx, query.GetAllForMonitoring)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]domain.MonitoringLocation, 0)
	for rows.Next() {
		locationModel := model.TrackedLocation{}
		var maxChatID int64
		if err := rows.Scan(
			&locationModel.ID,
			&locationModel.UserID,
			&locationModel.BusinessTypeID,
			&locationModel.Address,
			&locationModel.Location,
			&maxChatID,
			&locationModel.CreatedAt,
			&locationModel.UpdatedAt,
		); err != nil {
			return nil, err
		}

		location := dto.ToDomain(locationModel)
		locations = append(locations, *domain.NewMonitoringLocation(location, maxChatID))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}
