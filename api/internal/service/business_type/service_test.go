package businesstype

import (
	"context"
	"errors"
	"io"
	"testing"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	domainerrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

func TestServiceUpsert(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	service := NewService(testLogger(), repo)
	infraTypeID := uuid.New()

	result, err := service.Upsert(context.Background(), infraTypeID)
	if err != nil {
		t.Fatalf("Upsert returned an error: %v", err)
	}
	if result.InfraTypeID != infraTypeID {
		t.Fatalf("infra type ID: got %s, want %s", result.InfraTypeID, infraTypeID)
	}
	if repo.upsertCalls != 1 {
		t.Fatalf("repository upsert calls: got %d, want 1", repo.upsertCalls)
	}
}

func TestServiceUpsertMapsMissingInfraType(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{upsertErr: domainerrs.ErrInfraTypeNotFound}
	service := NewService(testLogger(), repo)

	_, err := service.Upsert(context.Background(), uuid.New())
	assertAppError(t, err, domainerrs.ErrInfraTypeNotFound)
}

func TestServiceGetByIDMapsNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{getByIDErr: domainerrs.ErrBusinessTypeNotFound}
	service := NewService(testLogger(), repo)

	_, err := service.GetByID(context.Background(), uuid.New())
	assertAppError(t, err, domainerrs.ErrBusinessTypeNotFound)
}

type fakeRepository struct {
	upsertCalls int
	upsertErr   error
	getByIDErr  error
}

func (r *fakeRepository) Upsert(
	_ context.Context,
	businessType *domain.BusinessType,
) (*domain.BusinessType, error) {
	r.upsertCalls++
	if r.upsertErr != nil {
		return nil, r.upsertErr
	}
	return businessType, nil
}

func (r *fakeRepository) GetByID(context.Context, uuid.UUID) (*domain.BusinessType, error) {
	return nil, r.getByIDErr
}

func (r *fakeRepository) GetAll(context.Context) ([]domain.BusinessType, error) {
	return nil, nil
}

func (r *fakeRepository) Delete(context.Context, uuid.UUID) error {
	return nil
}

func assertAppError(t *testing.T, err error, cause error) {
	t.Helper()

	var appErr *apperrs.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type: got %T, want *app.Error", err)
	}
	if appErr.Code != "not_found" {
		t.Fatalf("error code: got %q, want not_found", appErr.Code)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("error does not preserve cause %v", cause)
	}
}

func testLogger() *logger.Logger {
	return logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
}
