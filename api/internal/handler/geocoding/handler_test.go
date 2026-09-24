package geocoding

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

func TestSuggestAppliesLimit(t *testing.T) {
	t.Parallel()

	handler := New(&fakeGeocodingService{})
	request := httptest.NewRequest(http.MethodGet, "/api/v1/geocoding/suggestions?query=test&limit=1", nil)
	response := httptest.NewRecorder()

	handler.Suggest(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); got != "[{\"address\":\"first\",\"lat\":1,\"lon\":2}]\n" {
		t.Fatalf("body: got %q", got)
	}
}

func TestAddressPassesCoordinatesToService(t *testing.T) {
	t.Parallel()

	service := &fakeGeocodingService{}
	handler := New(service)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/geocoding/address?lat=59.94&lon=30.32",
		nil,
	)
	response := httptest.NewRecorder()

	handler.Address(response, request)

	if service.lat != 59.94 || service.lon != 30.32 {
		t.Fatalf(
			"service coordinates: got lat=%v lon=%v, want lat=59.94 lon=30.32",
			service.lat,
			service.lon,
		)
	}
}

type fakeGeocodingService struct {
	lat float64
	lon float64
}

func (*fakeGeocodingService) Suggest(context.Context, string) ([]domain.Address, error) {
	return []domain.Address{
		{Address: "first", Lat: 1, Lon: 2},
		{Address: "second", Lat: 3, Lon: 4},
	}, nil
}

func (s *fakeGeocodingService) Reverse(
	_ context.Context,
	geopoint *domain.GeoPoint,
) (*domain.Address, error) {
	s.lat = geopoint.Lat
	s.lon = geopoint.Lng
	return &domain.Address{Address: "first", Lat: geopoint.Lat, Lon: geopoint.Lng}, nil
}
