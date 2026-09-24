package geocoder

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAutocomplete(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/suggest/address" {
			t.Errorf("request = %s %s, want POST /suggest/address", r.Method, r.URL.Path)
		}
		assertHeaders(t, r)
		var request suggestRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.Query != "Невс 6" || request.Count != 5 || !request.RestrictValue {
			t.Errorf("unexpected request: %+v", request)
		}
		if len(request.Locations) != 1 || request.Locations[0].City != "Санкт-Петербург" {
			t.Errorf("unexpected locations: %+v", request.Locations)
		}
		return jsonResponse(http.StatusOK, `{"suggestions":[{"value":"Невский проспект, 6","data":{"geo_lat":"59.94","geo_lon":"30.32"}}]}`), nil
	})}

	result, err := New(httpClient, "https://suggestions.example/", "token").Autocomplete(context.Background(), "Невс 6")
	if err != nil {
		t.Fatalf("Autocomplete returned an error: %v", err)
	}
	if len(result) != 1 || result[0].Name != "Невский проспект, 6" || result[0].Latitude != 59.94 || result[0].Longitude != 30.32 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestReverse(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/geolocate/address" {
			t.Errorf("request = %s %s, want POST /geolocate/address", r.Method, r.URL.Path)
		}
		assertHeaders(t, r)
		var request geolocateRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.Lat != 59.94 || request.Lon != 30.32 || request.Count != 1 {
			t.Errorf("unexpected request: %+v", request)
		}
		return jsonResponse(http.StatusOK, `{"suggestions":[{"value":"Невский проспект, 6","data":{"geo_lat":"59.94","geo_lon":"30.32"}}]}`), nil
	})}

	result, err := New(httpClient, "https://suggestions.example", "token").Reverse(context.Background(), 30.32, 59.94)
	if err != nil {
		t.Fatalf("Reverse returned an error: %v", err)
	}
	if result.Address != "Невский проспект, 6" || result.Center.X != 30.32 || result.Center.Y != 59.94 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestReverseReturnsErrorForEmptyResponse(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"suggestions":[]}`), nil
	})}

	if _, err := New(httpClient, "https://suggestions.example", "token").Reverse(context.Background(), 30.32, 59.94); err == nil {
		t.Fatal("Reverse returned nil error")
	}
}

func TestClientReturnsErrorForNonOKStatus(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusUnauthorized, `{}`), nil
	})}

	if _, err := New(httpClient, "https://suggestions.example", "token").Autocomplete(context.Background(), "Невс 6"); err == nil {
		t.Fatal("Autocomplete returned nil error")
	}
}

func assertHeaders(t *testing.T, request *http.Request) {
	t.Helper()
	if got := request.Header.Get("Authorization"); got != "Token token" {
		t.Errorf("Authorization = %q, want %q", got, "Token token")
	}
	if got := request.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
