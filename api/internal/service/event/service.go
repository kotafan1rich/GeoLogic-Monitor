package event

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	domainerrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

const (
	defaultRadius uint16 = 500
	maxRadius     uint16 = 5000
)

type EventRepository interface {
	Upsert(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetUnnotifiedByPeriod(ctx context.Context, from, to time.Time) ([]*domain.Event, error)
	GetUnnotifiedNear(
		ctx context.Context,
		geoPoint *domain.GeoPoint,
		radius uint16,
		from *time.Time,
		to *time.Time,
	) ([]*domain.Event, error)
	MarkNotified(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo EventRepository
	log  *logger.Logger
}

func NewEventService(log *logger.Logger, repo EventRepository) *service {
	return &service{repo: repo, log: log}
}

func (s *service) Upsert(
	ctx context.Context,
	provider string,
	externalID string,
	lat float64,
	lon float64,
	date time.Time,
	info *string,
) (*domain.Event, error) {
	geoPoint, err := domain.NewGeoPoint(lat, lon)
	if err != nil {
		return nil, apperrs.ValidationError(err)
	}

	event, err := domain.NewEvent(provider, externalID, geoPoint, date, info)
	if err != nil {
		return nil, apperrs.ValidationError(err)
	}

	event, err = s.repo.Upsert(ctx, event)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to upsert event",
			slog.String("provider", provider),
			slog.String("external_id", externalID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return event, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, s.mapIDError(ctx, err, "failed to get event", id)
	}

	return event, nil
}

func (s *service) GetUnnotifiedByPeriod(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]*domain.Event, error) {
	events, err := s.repo.GetUnnotifiedByPeriod(ctx, from, to)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get unnotified events by period",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return events, nil
}

func (s *service) GetUnnotifiedNear(
	ctx context.Context,
	geoPoint *domain.GeoPoint,
	radius *uint16,
	from, to *time.Time,
) ([]*domain.Event, error) {
	finalRadius := defaultRadius
	if radius != nil {
		finalRadius = *radius
		if finalRadius == 0 || finalRadius > maxRadius {
			return nil, apperrs.ValidationError(domainerrs.ErrInvalidRadius)
		}
	}

	events, err := s.repo.GetUnnotifiedNear(ctx, geoPoint, finalRadius, from, to)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get nearby unnotified events",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return events, nil
}

func (s *service) MarkNotified(ctx context.Context, id uuid.UUID) error {
	err := s.repo.MarkNotified(ctx, id)
	if err != nil {
		return s.mapIDError(ctx, err, "failed to mark event as notified", id)
	}

	return nil
}

func (s *service) mapIDError(ctx context.Context, err error, message string, id uuid.UUID) error {
	if errors.Is(err, domainerrs.ErrEventNotFound) {
		return apperrs.Wrap(err, apperrs.ErrNotFound)
	}

	s.log.ErrorContext(ctx,
		message,
		slog.String("id", id.String()),
		slog.String("error", err.Error()),
	)
	return err
}
