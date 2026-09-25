package trackedlocation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
		&fakeRatingHistoryService{},
	)
	body := `{"name":"Кофейня на Невском","business_type_id":"` + businessTypeID.String() +
		`","address":"address","lat":59.93,"lon":30.32}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/tracked-locations", strings.NewReader(body))
	request.Header.Set("X-Max-User-Id", "42")
	response := httptest.NewRecorder()

	middleware.MaxUserID(http.HandlerFunc(handler.Create)).ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusCreated)
	}
	if locationService.createdUserID != userID ||
		locationService.createdBusinessTypeID != businessTypeID ||
		locationService.createdName != "Кофейня на Невском" {
		t.Fatalf(
			"unexpected service arguments: user=%s businessType=%s name=%q",
			locationService.createdUserID,
			locationService.createdBusinessTypeID,
			locationService.createdName,
		)
	}
}

func TestGetAllForMonitoringReturnsEmptyArray(t *testing.T) {
	t.Parallel()

	handler := New(fakeUserService{}, &fakeTrackedLocationService{}, &fakeRatingHistoryService{})
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

func TestGetRatingHistoryUsesDefaultMonths(t *testing.T) {
	t.Parallel()

	locationID := uuid.New()
	calculatedAt := time.Date(2026, time.September, 20, 10, 0, 0, 0, time.UTC)
	historyService := &fakeRatingHistoryService{
		history: &domain.LocationRatingHistory{
			TrackedLocationID: locationID,
			History: []*domain.CalculatedRating{
				{Value: 8, CalculatedAt: calculatedAt},
			},
		},
	}
	handler := New(fakeUserService{}, &fakeTrackedLocationService{}, historyService)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tracked-locations/"+locationID.String()+"/rating-history",
		nil,
	)
	request.SetPathValue("id", locationID.String())
	request.Header.Set("X-Max-User-Id", "42")
	response := httptest.NewRecorder()

	middleware.MaxUserID(http.HandlerFunc(handler.GetRatingHistory)).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusOK)
	}
	if historyService.maxUserID != 42 || historyService.locationID != locationID || historyService.months != 3 {
		t.Fatalf(
			"service arguments: got maxUserID=%d locationID=%s months=%d",
			historyService.maxUserID,
			historyService.locationID,
			historyService.months,
		)
	}
	wantBody := "[{\"value\":8,\"calculated_at\":\"2026-09-20T10:00:00Z\"}]\n"
	if got := response.Body.String(); got != wantBody {
		t.Fatalf("body: got %q, want %q", got, wantBody)
	}
}

func TestGetRatingHistoryRejectsInvalidMonths(t *testing.T) {
	t.Parallel()

	historyService := &fakeRatingHistoryService{}
	handler := New(fakeUserService{}, &fakeTrackedLocationService{}, historyService)
	locationID := uuid.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tracked-locations/"+locationID.String()+"/rating-history?months=0",
		nil,
	)
	request.SetPathValue("id", locationID.String())
	request.Header.Set("X-Max-User-Id", "42")
	response := httptest.NewRecorder()

	middleware.MaxUserID(http.HandlerFunc(handler.GetRatingHistory)).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusBadRequest)
	}
	if historyService.calls != 0 {
		t.Fatalf("service calls: got %d, want 0", historyService.calls)
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
	createdName           string
}

type fakeRatingHistoryService struct {
	history    *domain.LocationRatingHistory
	maxUserID  int64
	locationID uuid.UUID
	months     uint
	calls      int
}

func (s *fakeRatingHistoryService) GetHistory(
	_ context.Context,
	maxUserID int64,
	locationID uuid.UUID,
	months uint,
) (*domain.LocationRatingHistory, error) {
	s.calls++
	s.maxUserID = maxUserID
	s.locationID = locationID
	s.months = months
	if s.history == nil {
		return &domain.LocationRatingHistory{History: make([]*domain.CalculatedRating, 0)}, nil
	}
	return s.history, nil
}

func (s *fakeTrackedLocationService) Create(
	_ context.Context,
	userID uuid.UUID,
	businessTypeID uuid.UUID,
	name string,
	address string,
	lat float64,
	lon float64,
) (*domain.TrackedLocationRating, error) {
	s.createdUserID = userID
	s.createdBusinessTypeID = businessTypeID
	s.createdName = name
	return domain.NewTrackedLocationRating(
		&domain.TrackedLocation{
			ID:             uuid.New(),
			UserID:         userID,
			BusinessTypeID: businessTypeID,
			Name:           name,
			Address:        address,
			GeoPoint:       domain.GeoPoint{Lat: lat, Lon: lon},
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
