package geocoder

import (
	"context"
	"errors"
	"testing"

	geocoderintegration "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/geocoder"
)

func TestSuggestions(t *testing.T) {
	t.Parallel()

	client := &fakeClient{result: []geocoderintegration.Autocomplete{{
		Name:         "Невский проспект",
		BuildingName: "6",
		Latitude:     59.94,
		Longitude:    30.32,
	}}}
	repository := New(client)

	result, err := repository.Suggestions(context.Background(), "Невс 6")
	if err != nil {
		t.Fatalf("Suggestions returned an error: %v", err)
	}
	if client.query != "Невс 6" {
		t.Fatalf("client query: got %q, want %q", client.query, "Невс 6")
	}
	if len(result) != 1 {
		t.Fatalf("suggestion count: got %d, want 1", len(result))
	}
	if result[0].Address != "Невский проспект, 6" || result[0].Lat != 59.94 || result[0].Lon != 30.32 {
		t.Fatalf("unexpected address: %+v", result[0])
	}
}

func TestSuggestionsReturnsEmptySliceForEmptyResponse(t *testing.T) {
	t.Parallel()

	result, err := New(&fakeClient{}).Suggestions(context.Background(), "unknown")
	if err != nil {
		t.Fatalf("Suggestions returned an error: %v", err)
	}
	if result == nil || len(result) != 0 {
		t.Fatalf("result: got %#v, want empty non-nil slice", result)
	}
}

func TestSuggestionsPreservesClientError(t *testing.T) {
	t.Parallel()

	clientErr := errors.New("geocoder unavailable")
	_, err := New(&fakeClient{err: clientErr}).Suggestions(context.Background(), "Невс 6")
	if !errors.Is(err, clientErr) {
		t.Fatalf("error: got %v, want %v", err, clientErr)
	}
}

type fakeClient struct {
	result []geocoderintegration.Autocomplete
	err    error
	query  string
}

func (c *fakeClient) Autocomplete(
	_ context.Context,
	search string,
) ([]geocoderintegration.Autocomplete, error) {
	c.query = search
	return c.result, c.err
}
