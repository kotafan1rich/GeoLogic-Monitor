package calculate

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

func TestFormulaCalculator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		features   []domain.InfraTypeFeatures
		wantRating float64
	}{
		{name: "no infrastructure", wantRating: 5},
		{
			name:       "positive infrastructure",
			features:   []domain.InfraTypeFeatures{{Weight: 2, Influence: 1}},
			wantRating: 8,
		},
		{
			name:       "competitor",
			features:   []domain.InfraTypeFeatures{{Weight: 2, Influence: 1, IsCompetitor: true}},
			wantRating: 2,
		},
		{
			name: "equal positive and competitor",
			features: []domain.InfraTypeFeatures{
				{Weight: 1, Influence: 2},
				{Weight: 1, Influence: 2, IsCompetitor: true},
			},
			wantRating: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rating, err := NewFormulaCalculator().Calculate(context.Background(), domain.LocationFeatures{
				InfraTypes: tt.features,
			})
			if err != nil {
				t.Fatalf("Calculate returned an error: %v", err)
			}
			if rating.Value != tt.wantRating {
				t.Fatalf("rating: got %v, want %v", rating.Value, tt.wantRating)
			}
		})
	}
}

func TestFormulaCalculatorRatingIsFiniteAndBounded(t *testing.T) {
	t.Parallel()

	for _, influence := range []float64{0, 1, 1e100, math.MaxFloat64} {
		rating, err := NewFormulaCalculator().Calculate(context.Background(), domain.LocationFeatures{
			InfraTypes: []domain.InfraTypeFeatures{{Weight: 1, Influence: influence}},
		})
		if err != nil {
			t.Fatalf("Calculate returned an error: %v", err)
		}
		if math.IsNaN(rating.Value) || math.IsInf(rating.Value, 0) {
			t.Fatalf("rating is not finite: %v", rating.Value)
		}
		if rating.Value <= 0 || rating.Value >= 10 {
			t.Fatalf("rating is outside bounds: %v", rating.Value)
		}
		if rating.CalculatedAt.IsZero() {
			t.Fatal("calculation time is zero")
		}
		if rating.CalculatedAt.Location() != time.UTC {
			t.Fatalf("calculation time location: got %v, want UTC", rating.CalculatedAt.Location())
		}
	}
}

func TestFormulaCalculatorRespectsCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewFormulaCalculator().Calculate(ctx, domain.LocationFeatures{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error: got %v, want context.Canceled", err)
	}
}
