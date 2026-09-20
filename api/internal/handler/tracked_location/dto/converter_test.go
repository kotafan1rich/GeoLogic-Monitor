package dto

import (
	"testing"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

func TestToCreatedResponseIncludesCalculationTime(t *testing.T) {
	t.Parallel()

	calculatedAt := time.Date(2026, time.September, 20, 10, 0, 0, 0, time.UTC)
	response := ToCreatedResponse(domain.TrackedLocationRating{
		CalculatedRating: domain.CalculatedRating{
			Value:        8,
			CalculatedAt: calculatedAt,
		},
	})

	if response.Rating != 8 {
		t.Fatalf("rating: got %v, want 8", response.Rating)
	}
	if !response.RatingCalculatedAt.Equal(calculatedAt) {
		t.Fatalf("calculation time: got %v, want %v", response.RatingCalculatedAt, calculatedAt)
	}
}

func TestToResponseListIncludesNullableLatestRating(t *testing.T) {
	t.Parallel()

	calculatedAt := time.Date(2026, time.September, 20, 10, 0, 0, 0, time.UTC)
	responses := ToResponseList([]domain.TrackedLocation{
		{LatestRating: &domain.CalculatedRating{Value: 8, CalculatedAt: calculatedAt}},
		{},
	})

	if responses[0].Rating == nil || *responses[0].Rating != 8 {
		t.Fatalf("latest rating: got %v, want 8", responses[0].Rating)
	}
	if responses[0].RatingCalculatedAt == nil || !responses[0].RatingCalculatedAt.Equal(calculatedAt) {
		t.Fatalf("latest calculation time: got %v, want %v", responses[0].RatingCalculatedAt, calculatedAt)
	}
	if responses[1].Rating != nil || responses[1].RatingCalculatedAt != nil {
		t.Fatalf(
			"missing rating: got rating=%v calculated_at=%v, want nil values",
			responses[1].Rating,
			responses[1].RatingCalculatedAt,
		)
	}
}
