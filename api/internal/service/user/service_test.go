package user

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	domainerrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

func TestServiceUpsert(t *testing.T) {
	t.Parallel()

	repo := &fakeUserRepository{}
	service := NewUserService(testLogger(), repo)

	result, err := service.Upsert(context.Background(), 42, 84)
	if err != nil {
		t.Fatalf("Upsert returned an error: %v", err)
	}
	if result.MaxUserID != 42 || result.MaxChatID != 84 {
		t.Fatalf("unexpected user: %+v", result)
	}
	if repo.upsertCalls != 1 {
		t.Fatalf("repository upsert calls: got %d, want 1", repo.upsertCalls)
	}
}

func TestServiceUpsertRejectsInvalidMaxUserID(t *testing.T) {
	t.Parallel()

	repo := &fakeUserRepository{}
	service := NewUserService(testLogger(), repo)

	_, err := service.Upsert(context.Background(), 0, 84)
	var appErr *apperrs.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type: got %T, want *app.Error", err)
	}
	if appErr.Code != "validation_error" {
		t.Fatalf("error code: got %q, want validation_error", appErr.Code)
	}
	if !errors.Is(err, domainerrs.ErrInvalidMaxUserID) {
		t.Fatal("domain validation error is not preserved")
	}
	if repo.upsertCalls != 0 {
		t.Fatalf("repository upsert calls: got %d, want 0", repo.upsertCalls)
	}
}

func TestServiceGetByMaxUserIDMapsNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeUserRepository{getErr: domainerrs.ErrUserNotFound}
	service := NewUserService(testLogger(), repo)

	_, err := service.GetByMaxUserID(context.Background(), 42)
	var appErr *apperrs.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type: got %T, want *app.Error", err)
	}
	if appErr.Code != "not_found" {
		t.Fatalf("error code: got %q, want not_found", appErr.Code)
	}
	if !errors.Is(err, domainerrs.ErrUserNotFound) {
		t.Fatal("repository error is not preserved")
	}
}

type fakeUserRepository struct {
	upsertCalls int
	getErr      error
}

func (r *fakeUserRepository) Upsert(_ context.Context, user *domain.User) (*domain.User, error) {
	r.upsertCalls++
	return user, nil
}

func (r *fakeUserRepository) GetByMaxUserID(context.Context, int64) (*domain.User, error) {
	return nil, r.getErr
}

func testLogger() *logger.Logger {
	return logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
}
