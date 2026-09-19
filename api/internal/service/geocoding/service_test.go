package geocoding

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/geocoder"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

func TestSuggest(t *testing.T) {
	t.Parallel()

	client := &fakeClient{result: &geocoder.Autocomplete{
		Name:         "Невский проспект",
		BuildingName: "6",
		Latitude:     59.94,
		Longitude:    30.32,
	}}
	service := NewService(client, testLogger())

	result, err := service.Suggest(context.Background(), "  Невс 6  ")
	if err != nil {
		t.Fatalf("Suggest returned an error: %v", err)
	}
	if client.query != "Невс 6" {
		t.Fatalf("client query: got %q, want %q", client.query, "Невс 6")
	}
	if len(result) != 1 {
		t.Fatalf("suggestion count: got %d, want 1", len(result))
	}
	if result[0].Address != "Невский проспект, 6" || result[0].Lat != 59.94 || result[0].Lon != 30.32 {
		t.Fatalf("unexpected suggestion: %+v", result[0])
	}
}

func TestSuggestValidatesQuery(t *testing.T) {
	t.Parallel()

	client := &fakeClient{}
	service := NewService(client, testLogger())

	for _, query := range []string{" ", strings.Repeat("я", maxQueryLength+1)} {
		_, err := service.Suggest(context.Background(), query)
		var appErr *apperrs.Error
		if !errors.As(err, &appErr) || appErr.Code != "validation_error" {
			t.Fatalf("error for query length %d: got %v, want validation_error", len(query), err)
		}
	}
	if client.calls != 0 {
		t.Fatalf("client calls: got %d, want 0", client.calls)
	}
}

func TestSuggestMapsClientError(t *testing.T) {
	t.Parallel()

	clientErr := errors.New("geocoder unavailable")
	service := NewService(&fakeClient{err: clientErr}, testLogger())

	_, err := service.Suggest(context.Background(), "Невс 6")
	var appErr *apperrs.Error
	if !errors.As(err, &appErr) || appErr.Code != "provider_unavailable" {
		t.Fatalf("error: got %v, want provider_unavailable", err)
	}
	if !errors.Is(err, clientErr) {
		t.Fatal("client error is not preserved")
	}
}

type fakeClient struct {
	result *geocoder.Autocomplete
	err    error
	query  string
	calls  int
}

func (c *fakeClient) Autocomplete(_ context.Context, search string) (*geocoder.Autocomplete, error) {
	c.calls++
	c.query = search
	return c.result, c.err
}

func testLogger() *logger.Logger {
	return logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
}
