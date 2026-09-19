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

	handler := New(fakeGeocodingService{})
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

type fakeGeocodingService struct{}

func (fakeGeocodingService) Suggest(context.Context, string) ([]domain.Address, error) {
	return []domain.Address{
		{Address: "first", Lat: 1, Lon: 2},
		{Address: "second", Lat: 3, Lon: 4},
	}, nil
}
