package trackedlocation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/middleware"
)

func TestCreateUsesCurrentUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	businessTypeID := uuid.New()
	locationService := &fakeTrackedLocationService{}
	handler := New(
		fakeUserService{user: &domain.User{ID: userID, MaxUserID: 42}},
		locationService,
	)
	body := `{"business_type_id":"` + businessTypeID.String() +
		`","address":"address","lat":59.93,"lon":30.32}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/tracked-locations", strings.NewReader(body))
	request.Header.Set("X-Max-User-Id", "42")
	response := httptest.NewRecorder()

	middleware.MaxUserID(http.HandlerFunc(handler.Create)).ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusCreated)
	}
	if locationService.createdUserID != userID || locationService.createdBusinessTypeID != businessTypeID {
		t.Fatalf(
			"unexpected service IDs: user=%s businessType=%s",
			locationService.createdUserID,
			locationService.createdBusinessTypeID,
		)
	}
}

func TestGetAllForMonitoringReturnsEmptyArray(t *testing.T) {
	t.Parallel()

	handler := New(fakeUserService{}, &fakeTrackedLocationService{})
	request := httptest.NewRequest(http.MethodGet, "/internal/v1/tracked-locations", nil)
	response := httptest.NewRecorder()

	handler.GetAllForMonitoring(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); got != "[]\n" {
		t.Fatalf("body: got %q, want empty JSON array", got)
	}
}

type fakeUserService struct {
	user *domain.User
}

func (s fakeUserService) GetByMaxUserID(context.Context, int64) (*domain.User, error) {
	return s.user, nil
}

type fakeTrackedLocationService struct {
	createdUserID         uuid.UUID
	createdBusinessTypeID uuid.UUID
}

func (s *fakeTrackedLocationService) Create(
	_ context.Context,
	userID uuid.UUID,
	businessTypeID uuid.UUID,
	address string,
	lat float64,
	lng float64,
) (*domain.TrackedLocationRating, error) {
	s.createdUserID = userID
	s.createdBusinessTypeID = businessTypeID
	return domain.NewTrackedLocationRating(
		&domain.TrackedLocation{
			ID:             uuid.New(),
			UserID:         userID,
			BusinessTypeID: businessTypeID,
			Address:        address,
			GeoPoint:       domain.GeoPoint{Lat: lat, Lng: lng},
		},
		&domain.CalculatedRating{Value: 5},
	), nil
}

func (*fakeTrackedLocationService) GetByUserID(
	context.Context,
	uuid.UUID,
) ([]domain.TrackedLocation, error) {
	return []domain.TrackedLocation{}, nil
}

func (*fakeTrackedLocationService) DeleteForUser(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}

func (*fakeTrackedLocationService) GetAllForMonitoring(
	context.Context,
) ([]domain.MonitoringLocation, error) {
	return []domain.MonitoringLocation{}, nil
}
