package infra

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
	service := NewTypeService(testLogger(), repo)

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
	service := NewTypeService(testLogger(), repo)

	_, err := service.Upsert(context.Background(), "", "School", 2, 1000)
	assertAppError(t, err, "validation_error", domainerrs.ErrInvalidSlug)
	if repo.upsertCalls != 0 {
		t.Fatalf("repository upsert calls: got %d, want 0", repo.upsertCalls)
	}
}

func TestInfraServiceUpsert(t *testing.T) {
	t.Parallel()

	repo := &fakeInfraRepository{}
	service := NewInfraService(testLogger(), repo)
	typeID := uuid.New()

	result, err := service.Upsert(context.Background(), "external-id", typeID, 59.93, 30.32, "address", nil)
	if err != nil {
		t.Fatalf("Upsert returned an error: %v", err)
	}
	if result.ExternalID != "external-id" || result.TypeID != typeID || result.Address != "address" {
		t.Fatalf("unexpected infra object: %+v", result)
	}
	if repo.upsertCalls != 1 {
		t.Fatalf("repository upsert calls: got %d, want 1", repo.upsertCalls)
	}
}

func TestInfraServiceUpsertRejectsInvalidExternalID(t *testing.T) {
	t.Parallel()

	repo := &fakeInfraRepository{}
	service := NewInfraService(testLogger(), repo)

	_, err := service.Upsert(context.Background(), "", uuid.New(), 59.93, 30.32, "address", nil)
	assertAppError(t, err, "validation_error", domainerrs.ErrInvalidExternalID)
	if repo.upsertCalls != 0 {
		t.Fatalf("repository upsert calls: got %d, want 0", repo.upsertCalls)
	}
}

func TestInfraServiceUpsertMapsMissingType(t *testing.T) {
	t.Parallel()

	repo := &fakeInfraRepository{upsertErr: domainerrs.ErrInfraTypeNotFound}
	service := NewInfraService(testLogger(), repo)

	_, err := service.Upsert(
		context.Background(),
		"external-id",
		uuid.New(),
		59.93,
		30.32,
		"address",
		nil,
	)
	assertAppError(t, err, "not_found", domainerrs.ErrInfraTypeNotFound)
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

type fakeInfraRepository struct {
	upsertCalls int
	upsertErr   error
}

func (r *fakeInfraRepository) Upsert(
	_ context.Context,
	infraObject *domain.InfraObject,
) (*domain.InfraObject, error) {
	r.upsertCalls++
	if r.upsertErr != nil {
		return nil, r.upsertErr
	}
	return infraObject, nil
}

func (r *fakeInfraRepository) GetByID(context.Context, uuid.UUID) (*domain.InfraObject, error) {
	return nil, nil
}

func (r *fakeInfraRepository) Near(context.Context, *domain.GeoPoint) ([]*domain.InfraObject, error) {
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
