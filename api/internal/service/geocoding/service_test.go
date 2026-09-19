package geocoding

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

func TestSuggest(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{result: []domain.Address{{
		Address: "Невский проспект, 6",
		Lat:     59.94,
		Lon:     30.32,
	}}}
	service := NewService(repo, testLogger())

	result, err := service.Suggest(context.Background(), "  Невс 6  ")
	if err != nil {
		t.Fatalf("Suggest returned an error: %v", err)
	}
	if repo.query != "Невс 6" {
		t.Fatalf("repository query: got %q, want %q", repo.query, "Невс 6")
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

	repo := &fakeRepository{}
	service := NewService(repo, testLogger())

	for _, query := range []string{" ", strings.Repeat("я", maxQueryLength+1)} {
		_, err := service.Suggest(context.Background(), query)
		var appErr *apperrs.Error
		if !errors.As(err, &appErr) || appErr.Code != "validation_error" {
			t.Fatalf("error for query length %d: got %v, want validation_error", len(query), err)
		}
	}
	if repo.calls != 0 {
		t.Fatalf("repository calls: got %d, want 0", repo.calls)
	}
}

func TestSuggestMapsClientError(t *testing.T) {
	t.Parallel()

	repositoryErr := errors.New("geocoder unavailable")
	service := NewService(&fakeRepository{err: repositoryErr}, testLogger())

	_, err := service.Suggest(context.Background(), "Невс 6")
	var appErr *apperrs.Error
	if !errors.As(err, &appErr) || appErr.Code != "provider_unavailable" {
		t.Fatalf("error: got %v, want provider_unavailable", err)
	}
	if !errors.Is(err, repositoryErr) {
		t.Fatal("repository error is not preserved")
	}
}

type fakeRepository struct {
	result []domain.Address
	err    error
	query  string
	calls  int
}

func (r *fakeRepository) Suggestions(_ context.Context, query string) ([]domain.Address, error) {
	r.calls++
	r.query = query
	return r.result, r.err
}

func testLogger() *logger.Logger {
	return logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
}
