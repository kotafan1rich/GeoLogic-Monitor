package infratype

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

func TestTypeServiceUpsert(t *testing.T) {
	t.Parallel()

	repo := &fakeInfraTypeRepository{}
	service := NewService(testLogger(), repo)

	result, err := service.Upsert(context.Background(), "school", "School", 2, 1000)
	if err != nil {
		t.Fatalf("Upsert returned an error: %v", err)
	}
	if result.Slug != "school" || result.Name != "School" {
		t.Fatalf("unexpected infra type: %+v", result)
	}
	if repo.upsertCalls != 1 {
		t.Fatalf("repository upsert calls: got %d, want 1", repo.upsertCalls)
	}
}

func TestTypeServiceUpsertRejectsInvalidSlug(t *testing.T) {
	t.Parallel()

	repo := &fakeInfraTypeRepository{}
	service := NewService(testLogger(), repo)

	_, err := service.Upsert(context.Background(), "", "School", 2, 1000)
	assertAppError(t, err, "validation_error", domainerrs.ErrInvalidSlug)
	if repo.upsertCalls != 0 {
		t.Fatalf("repository upsert calls: got %d, want 0", repo.upsertCalls)
	}
}

type fakeInfraTypeRepository struct {
	upsertCalls int
}

func (r *fakeInfraTypeRepository) Upsert(
	_ context.Context,
	infraType *domain.InfraType,
) (*domain.InfraType, error) {
	r.upsertCalls++
	return infraType, nil
}

func (r *fakeInfraTypeRepository) GetByID(context.Context, uuid.UUID) (*domain.InfraType, error) {
	return nil, nil
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
