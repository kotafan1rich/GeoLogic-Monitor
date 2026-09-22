package event

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	domainerrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

func TestServiceUpsert(t *testing.T) {
	t.Parallel()

	repo := &fakeEventRepository{}
	service := NewEventService(testLogger(), repo)
	date := time.Now()

	result, err := service.Upsert(context.Background(), "provider", "external-id", 59.93, 30.32, date, nil)
	if err != nil {
		t.Fatalf("Upsert returned an error: %v", err)
	}
	if result.Provider != "provider" || result.ExternalID != "external-id" || !result.Date.Equal(date) {
		t.Fatalf("unexpected event: %+v", result)
	}
	if repo.upsertCalls != 1 {
		t.Fatalf("repository upsert calls: got %d, want 1", repo.upsertCalls)
	}
}

func TestServiceUpsertRejectsInvalidProvider(t *testing.T) {
	t.Parallel()

	repo := &fakeEventRepository{}
	service := NewEventService(testLogger(), repo)

	_, err := service.Upsert(context.Background(), "", "external-id", 59.93, 30.32, time.Now(), nil)
	assertAppError(t, err, "validation_error", domainerrs.ErrInvalidEventProvider)
	if repo.upsertCalls != 0 {
		t.Fatalf("repository upsert calls: got %d, want 0", repo.upsertCalls)
	}
}

func TestServiceGetByIDMapsNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeEventRepository{getByIDErr: domainerrs.ErrEventNotFound}
	service := NewEventService(testLogger(), repo)

	_, err := service.GetByID(context.Background(), uuid.New())
	assertAppError(t, err, "not_found", domainerrs.ErrEventNotFound)
}

func TestServiceGetUnnotifiedNearUsesDefaultRadius(t *testing.T) {
	t.Parallel()

	repo := &fakeEventRepository{}
	service := NewEventService(testLogger(), repo)
	geoPoint, err := domain.NewGeoPoint(59.93, 30.32)
	if err != nil {
		t.Fatalf("NewGeoPoint returned an error: %v", err)
	}

	_, err = service.GetUnnotifiedNear(context.Background(), geoPoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("GetUnnotifiedNear returned an error: %v", err)
	}
	if repo.nearRadius != defaultRadius {
		t.Fatalf("repository radius: got %d, want %d", repo.nearRadius, defaultRadius)
	}
}

func TestServiceGetUnnotifiedNearRejectsInvalidRadius(t *testing.T) {
	t.Parallel()

	repo := &fakeEventRepository{}
	service := NewEventService(testLogger(), repo)
	geoPoint, err := domain.NewGeoPoint(59.93, 30.32)
	if err != nil {
		t.Fatalf("NewGeoPoint returned an error: %v", err)
	}
	radius := uint16(0)

	_, err = service.GetUnnotifiedNear(context.Background(), geoPoint, &radius, nil, nil)
	assertAppError(t, err, "validation_error", domainerrs.ErrInvalidRadius)
	if repo.nearCalls != 0 {
		t.Fatalf("repository near calls: got %d, want 0", repo.nearCalls)
	}
}

type fakeEventRepository struct {
	upsertCalls int
	getByIDErr  error
	nearCalls   int
	nearRadius  uint16
}

func (r *fakeEventRepository) Upsert(
	_ context.Context,
	event *domain.Event,
) (*domain.Event, error) {
	r.upsertCalls++
	return event, nil
}

func (r *fakeEventRepository) GetByID(context.Context, uuid.UUID) (*domain.Event, error) {
	return nil, r.getByIDErr
}

func (r *fakeEventRepository) GetUnnotifiedByPeriod(
	context.Context,
	time.Time,
	time.Time,
) ([]*domain.Event, error) {
	return nil, nil
}

func (r *fakeEventRepository) GetUnnotifiedNear(
	_ context.Context,
	_ *domain.GeoPoint,
	radius uint16,
	_ *time.Time,
	_ *time.Time,
) ([]*domain.Event, error) {
	r.nearCalls++
	r.nearRadius = radius
	return nil, nil
}

func (r *fakeEventRepository) MarkNotified(context.Context, uuid.UUID) error {
	return nil
}

func assertAppError(t *testing.T, err error, code string, cause error) {
	t.Helper()

	var appErr *apperrs.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type: got %T, want *app.Error", err)
	}
	if appErr.Code != code {
		t.Fatalf("error code: got %q, want %q", appErr.Code, code)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("error does not preserve cause %v", cause)
	}
}

func testLogger() *logger.Logger {
	return logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
}
