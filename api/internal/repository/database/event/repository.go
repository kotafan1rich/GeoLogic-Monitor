package event

import (
	"context"
	"errors"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/event/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/event/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/event/query"
	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	eventModel := dto.ToModel(event)
	err := scanEvent(
		r.db.QueryRow(
			ctx,
			query.Upsert,
			eventModel.Provider,
			eventModel.ExternalID,
			eventModel.Location,
			eventModel.Date,
			eventModel.Info,
		),
		eventModel,
	)
	if err != nil {
		return nil, err
	}

	return dto.ToDomain(eventModel), nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	eventModel := model.Event{}
	err := scanEvent(r.db.QueryRow(ctx, query.GetByID, id), &eventModel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrEventNotFound
		}
		return nil, err
	}

	return dto.ToDomain(&eventModel), nil
}

func (r *repository) GetUnnotifiedByPeriod(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]*domain.Event, error) {
	return r.getMany(ctx, query.GetUnnotifiedByPeriod, from, to)
}

func (r *repository) GetUnnotifiedNear(
	ctx context.Context,
	geopoint *domain.GeoPoint,
	radius uint16,
	from *time.Time,
	to *time.Time,
) ([]*domain.Event, error) {
	location := basemodel.GeoPoint(*geopoint)
	return r.getMany(ctx, query.GetUnnotifiedNear, location, radius, from, to)
}

func (r *repository) MarkNotified(ctx context.Context, id uuid.UUID) error {
	var notifiedAt time.Time
	err := r.db.QueryRow(ctx, query.MarkNotified, id).Scan(&notifiedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrEventNotFound
		}
		return err
	}
	return nil
}

func (r *repository) getMany(
	ctx context.Context,
	queryString string,
	args ...any,
) ([]*domain.Event, error) {
	rows, err := r.db.Query(ctx, queryString, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*domain.Event, 0)
	for rows.Next() {
		eventModel := model.Event{}
		if err := scanEvent(rows, &eventModel); err != nil {
			return nil, err
		}
		events = append(events, dto.ToDomain(&eventModel))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanEvent(row scanner, event *model.Event) error {
	return row.Scan(
		&event.ID,
		&event.Provider,
		&event.ExternalID,
		&event.Location,
		&event.Date,
		&event.Info,
		&event.NotifiedAt,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
}
