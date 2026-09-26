package competitor

import (
	"context"
	"fmt"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/database/postgresql"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/repository/postgresql/geo"
)

type CompetitorRepo struct {
	db *postgresql.DB
}

func New(db *postgresql.DB) *CompetitorRepo {
	return &CompetitorRepo{db: db}
}

func (r *CompetitorRepo) Add(
	ctx context.Context,
	trackedLocationID string,
	external_id string,
	name string,
	typeID string,
	address string,
	location geo.GeoPoint,
	openedAt time.Time,
	ttl time.Duration,
) error {
	const op = "competitor.CompetitorRepo.Add"

	_, err := r.db.Conn(ctx).Exec(
		ctx,
		AddCompetitorQuery,
		trackedLocationID,
		external_id,
		name,
		typeID,
		address,
		location,
		openedAt,
		time.Now().Add(ttl),
	)
	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrInsertData, err)
	}

	return nil
}

func (r *CompetitorRepo) UpdateStatus(
	ctx context.Context,
	trackedLocationID string,
	external_id string,
	notifiedAt time.Time,
) error {
	const op = "competitor.CompetitorRepo.UpdateStatus"

	_, err := r.db.Conn(ctx).Exec(
		ctx,
		ChangeStatusQuery,
		trackedLocationID,
		external_id,
		notifiedAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrUpdateStatus, err)
	}

	return nil
}
