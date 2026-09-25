package route

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
)

func TestWalkingDistances(t *testing.T) {
	t.Parallel()

	distance := 1557.6
	service := &fakeOSRMService{distances: []*float64{&distance, nil}}
	handler := New(service)
	request := httptest.NewRequest(
		http.MethodPost,
		"/internal/v1/routes/walking-distances",
		strings.NewReader(`{
			"source":{"lat":59.939095,"lon":30.315868},
			"destinations":[
				{"lat":59.934280,"lon":30.335098},
				{"lat":59.936000,"lon":30.327000}
			]
		}`),
	)
	response := httptest.NewRecorder()

	handler.WalkingDistances(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); got != "{\"distances\":[1557.6,null]}\n" {
		t.Fatalf("body: got %q", got)
	}
	if service.src == nil || service.src.Lat != 59.939095 || service.src.Lon != 30.315868 {
		t.Fatalf("unexpected source: %+v", service.src)
	}
	if len(service.dst) != 2 || service.dst[1].Lat != 59.936 || service.dst[1].Lon != 30.327 {
		t.Fatalf("unexpected destinations: %+v", service.dst)
	}
}

func TestWalkingDistancesRejectsEmptyDestinations(t *testing.T) {
	t.Parallel()

	service := &fakeOSRMService{}
	handler := New(service)
	request := httptest.NewRequest(
		http.MethodPost,
		"/internal/v1/routes/walking-distances",
		strings.NewReader(`{"source":{"lat":59.93,"lon":30.32},"destinations":[]}`),
	)
	response := httptest.NewRecorder()

	handler.WalkingDistances(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusBadRequest)
	}
	if service.calls != 0 {
		t.Fatalf("service calls: got %d, want 0", service.calls)
	}
}

func TestWalkingDistancesReturnsProviderError(t *testing.T) {
	t.Parallel()

	service := &fakeOSRMService{
		err: app.Wrap(errors.New("osrm unavailable"), app.ErrProviderUnavailable),
	}
	handler := New(service)
	request := httptest.NewRequest(
		http.MethodPost,
		"/internal/v1/routes/walking-distances",
		strings.NewReader(`{
			"source":{"lat":59.93,"lon":30.32},
			"destinations":[{"lat":59.94,"lon":30.33}]
		}`),
	)
	response := httptest.NewRecorder()

	handler.WalkingDistances(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
}

type fakeOSRMService struct {
	distances []*float64
	err       error
	src       *domain.GeoPoint
	dst       []*domain.GeoPoint
	calls     int
}

func (s *fakeOSRMService) WalkingDistances(
	_ context.Context,
	src *domain.GeoPoint,
	dst []*domain.GeoPoint,
) ([]*float64, error) {
	s.calls++
	s.src = src
	s.dst = dst
	return s.distances, s.err
}
