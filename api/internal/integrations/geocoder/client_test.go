package geocoder

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestParseFree(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/parse/free" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/parse/free")
		}
		if got := r.URL.Query().Get("street"); got != "Невский 6" {
			t.Errorf("street = %q, want %q", got, "Невский 6")
		}
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("Authorization = %q, want an empty header", got)
		}

		return jsonResponse(http.StatusOK, `{"Building_ID":10,"ID":20,"Name":"Невский проспект, 6","Longitude":30.32,"Latitude":59.94,"District_ID":30,"District":"Центральный","Flat":""}`), nil
	})}

	client := New(httpClient, "https://geocoder.example/")
	address, err := client.ParseFree(context.Background(), "Невский 6")
	if err != nil {
		t.Fatalf("ParseFree returned an error: %v", err)
	}
	if address.BuildingID != 10 || address.ID != 20 || address.Name != "Невский проспект, 6" {
		t.Fatalf("unexpected address: %+v", address)
	}
}

func TestReverse(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/geocode/reverse" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/geocode/reverse")
		}
		if got := r.URL.Query().Get("x"); got != "30.32" {
			t.Errorf("x = %q, want %q", got, "30.32")
		}
		if got := r.URL.Query().Get("y"); got != "59.94" {
			t.Errorf("y = %q, want %q", got, "59.94")
		}

		return jsonResponse(http.StatusOK, `{"id":20,"address":"Невский проспект, 6","center":{"x":30.32,"y":59.94},"distance":5.25}`), nil
	})}

	client := New(httpClient, "https://geocoder.example")
	result, err := client.Reverse(context.Background(), 30.32, 59.94)
	if err != nil {
		t.Fatalf("Reverse returned an error: %v", err)
	}
	if result.Address != "Невский проспект, 6" || result.Center.X != 30.32 || result.Center.Y != 59.94 || result.Distance != 5.25 {
		t.Fatalf("unexpected geocode: %+v", result)
	}
}

func TestAutocomplete(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/autocomplete/universal" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/autocomplete/universal")
		}
		if got := r.URL.Query().Get("s"); got != "Невс 6" {
			t.Errorf("s = %q, want %q", got, "Невс 6")
		}

		return jsonResponse(http.StatusOK, `[{"address_id":20,"Name":"Невский проспект","building_id":10,"building_name":"6","Longitude":30.32,"Latitude":59.94}]`), nil
	})}

	client := New(httpClient, "http://geocoder.example")
	result, err := client.Autocomplete(context.Background(), "Невс 6")
	if err != nil {
		t.Fatalf("Autocomplete returned an error: %v", err)
	}
	if len(result) != 1 || result[0].Name != "Невский проспект" || result[0].BuildingName != "6" {
		t.Fatalf("unexpected autocomplete result: %+v", result)
	}
}

func TestClientReturnsErrorForNonOKStatus(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusUnprocessableEntity, `{}`), nil
	})}

	client := New(httpClient, "https://geocoder.example")
	_, err := client.Autocomplete(context.Background(), "Невс 6")
	if err == nil {
		t.Fatal("Autocomplete returned nil error")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
