package rating

import (
	"context"
	"time"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/pgerrors"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/rating/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/rating/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/rating/query"
)

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *repository {
	return &repository{db: db}
}

func (r *repository) Create(
	ctx context.Context,
	locationRating *domain.LocationRating,
) (*domain.LocationRating, error) {
	locationRatingModel := dto.ToModel(*locationRating)
	if locationRatingModel.ID == uuid.Nil() {
		locationRatingModel.ID = uuid.New()
	}

	err := scanLocationRating(
		r.db.QueryRow(
			ctx,
			query.Create,
			locationRatingModel.ID,
			locationRatingModel.TrackedLocationID,
			locationRatingModel.Value,
			locationRatingModel.CalculatedAt,
		),
		locationRatingModel,
	)
	if err != nil {
		return nil, translateWriteError(err)
	}

	return dto.ToDomain(*locationRatingModel), nil
}

func (r *repository) CreateMany(
	ctx context.Context,
	locationRatings []domain.LocationRating,
) ([]domain.LocationRating, error) {
	if len(locationRatings) == 0 {
		return []domain.LocationRating{}, nil
	}

	ids := make([]uuid.UUID, len(locationRatings))
	trackedLocationIDs := make([]uuid.UUID, len(locationRatings))
	values := make([]float64, len(locationRatings))
	calculatedAt := make([]time.Time, len(locationRatings))
	for i := range locationRatings {
		if locationRatings[i].ID == uuid.Nil() {
			locationRatings[i].ID = uuid.New()
		}
		ids[i] = locationRatings[i].ID
		trackedLocationIDs[i] = locationRatings[i].TrackedLocationID
		values[i] = locationRatings[i].Value
		calculatedAt[i] = locationRatings[i].CalculatedAt
	}

	rows, err := r.db.Query(
		ctx,
		query.CreateMany,
		ids,
		trackedLocationIDs,
		values,
		calculatedAt,
	)
	if err != nil {
		return nil, translateWriteError(err)
	}
	defer rows.Close()

	createdRatings := make([]domain.LocationRating, 0, len(locationRatings))
	for rows.Next() {
		locationRatingModel := model.LocationRating{}
		if err := scanLocationRating(rows, &locationRatingModel); err != nil {
			return nil, err
		}
		createdRatings = append(createdRatings, *dto.ToDomain(locationRatingModel))
	}

	if err := rows.Err(); err != nil {
		return nil, translateWriteError(err)
	}

	return createdRatings, nil
}

func (r *repository) GetHistory(
	ctx context.Context,
	trackedLocationID uuid.UUID,
	months uint,
) (*domain.LocationRatingHistory, error) {
	rows, err := r.db.Query(ctx, query.GetHistory, trackedLocationID, int(months))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locationRatings := make([]*model.LocationRating, 0)
	for rows.Next() {
		locationRatingModel := model.LocationRating{}
		if err := scanLocationRating(rows, &locationRatingModel); err != nil {
			return nil, err
		}
		locationRatings = append(locationRatings, &locationRatingModel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dto.ToDomainHistory(trackedLocationID, locationRatings), nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanLocationRating(row scanner, locationRating *model.LocationRating) error {
	return row.Scan(
		&locationRating.ID,
		&locationRating.TrackedLocationID,
		&locationRating.Value,
		&locationRating.CalculatedAt,
		&locationRating.CreatedAt,
		&locationRating.UpdatedAt,
	)
}

func translateWriteError(err error) error {
	if pgerrors.IsForeignKeyViolation(err) {
		return errs.ErrTrackedLocationNotFound
	}
	return err
}
