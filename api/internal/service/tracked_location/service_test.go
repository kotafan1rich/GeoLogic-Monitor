package trackedlocation

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

func TestServiceCreate(t *testing.T) {
	t.Parallel()

	repo := &fakeTrackedLocationRepository{}
	service := newTestService(repo)
	userID := uuid.New()
	businessTypeID := uuid.New()

	result, err := service.Create(context.Background(), userID, businessTypeID, "address", 59.93, 30.32)
	if err != nil {
		t.Fatalf("Create returned an error: %v", err)
	}
	if result.UserID != userID || result.BusinessTypeID != businessTypeID || result.Address != "address" {
		t.Fatalf("unexpected tracked location: %+v", result)
	}
	if result.Value != 5 {
		t.Fatalf("rating: got %v, want 5", result.Value)
	}
	if repo.createCalls != 1 {
		t.Fatalf("repository create calls: got %d, want 1", repo.createCalls)
	}
}

func TestServiceCreateRejectsInvalidLatitude(t *testing.T) {
	t.Parallel()

	repo := &fakeTrackedLocationRepository{}
	service := newTestService(repo)

	_, err := service.Create(context.Background(), uuid.New(), uuid.New(), "address", 91, 30.32)
	var appErr *apperrs.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type: got %T, want *app.Error", err)
	}
	if appErr.Code != "validation_error" {
		t.Fatalf("error code: got %q, want validation_error", appErr.Code)
	}
	if !errors.Is(err, domainerrs.ErrInvalidLat) {
		t.Fatal("domain validation error is not preserved")
	}
	if repo.createCalls != 0 {
		t.Fatalf("repository create calls: got %d, want 0", repo.createCalls)
	}
}

func TestServiceGetByIDMapsNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeTrackedLocationRepository{getByIDErr: domainerrs.ErrTrackedLocationNotFound}
	service := newTestService(repo)

	_, err := service.GetByID(context.Background(), uuid.New())
	var appErr *apperrs.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type: got %T, want *app.Error", err)
	}
	if appErr.Code != "not_found" {
		t.Fatalf("error code: got %q, want not_found", appErr.Code)
	}
	if !errors.Is(err, domainerrs.ErrTrackedLocationNotFound) {
		t.Fatal("repository error is not preserved")
	}
}

func TestServiceDeleteForUserHidesForeignLocation(t *testing.T) {
	t.Parallel()

	repo := &fakeTrackedLocationRepository{
		getByIDResult: &domain.TrackedLocation{UserID: uuid.New()},
	}
	service := newTestService(repo)

	err := service.DeleteForUser(context.Background(), uuid.New(), uuid.New())
	var appErr *apperrs.Error
	if !errors.As(err, &appErr) || appErr.Code != "not_found" {
		t.Fatalf("error: got %v, want not_found", err)
	}
	if repo.deleteCalls != 0 {
		t.Fatalf("repository delete calls: got %d, want 0", repo.deleteCalls)
	}
}

type fakeTrackedLocationRepository struct {
	createCalls   int
	deleteCalls   int
	getByIDResult *domain.TrackedLocation
	getByIDErr    error
}

type fakeOSRMRepository struct{}

func (*fakeOSRMRepository) FilterWalkingDistance(
	context.Context,
	*domain.GeoPoint,
	[]*domain.InfraObject,
) ([]*domain.InfraObjectDistance, error) {
	return []*domain.InfraObjectDistance{}, nil
}

type fakeInfraService struct{}

func (*fakeInfraService) Near(
	context.Context,
	*domain.GeoPoint,
) ([]*domain.InfraObject, error) {
	return []*domain.InfraObject{}, nil
}

type fakeBusinessTypeService struct{}

func (*fakeBusinessTypeService) GetByID(
	_ context.Context,
	id uuid.UUID,
) (*domain.BusinessType, error) {
	return &domain.BusinessType{ID: id, InfraTypeID: uuid.New()}, nil
}

type fakeRatingService struct{}

func (*fakeRatingService) Calculate(
	context.Context,
	domain.LocationFeatures,
) (*domain.CalculatedRating, error) {
	return &domain.CalculatedRating{Value: 5}, nil
}

func (*fakeRatingService) Create(
	context.Context,
	uuid.UUID,
	*domain.CalculatedRating,
) (*domain.LocationRating, error) {
	return &domain.LocationRating{}, nil
}

type fakeTxManager struct{}

func (fakeTxManager) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func newTestService(repo TrackedLocationRepository) *service {
	return NewTrackedLocationService(
		repo,
		&fakeOSRMRepository{},
		&fakeInfraService{},
		&fakeBusinessTypeService{},
		&fakeRatingService{},
		fakeTxManager{},
		testLogger(),
	)
}

func (r *fakeTrackedLocationRepository) Create(
	_ context.Context,
	location *domain.TrackedLocation,
) (*domain.TrackedLocation, error) {
	r.createCalls++
	return location, nil
}

func (r *fakeTrackedLocationRepository) GetByID(
	context.Context,
	uuid.UUID,
) (*domain.TrackedLocation, error) {
	return r.getByIDResult, r.getByIDErr
}

func (r *fakeTrackedLocationRepository) GetByUserID(
	context.Context,
	uuid.UUID,
) ([]domain.TrackedLocation, error) {
	return nil, nil
}

func (r *fakeTrackedLocationRepository) Delete(context.Context, uuid.UUID) error {
	r.deleteCalls++
	return nil
}

func (r *fakeTrackedLocationRepository) GetAllForMonitoring(
	context.Context,
) ([]domain.MonitoringLocation, error) {
	return nil, nil
}

func testLogger() *logger.Logger {
	return logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
}
