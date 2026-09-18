package rating

import (
	"context"
	"math"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

type FormulaCalculator struct{}

func NewFormulaCalculator() FormulaCalculator {
	return FormulaCalculator{}
}

func (FormulaCalculator) Calculate(
	ctx context.Context,
	location domain.LocationFeatures,
) (*domain.CalculatedRating, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var positive, competition float64
	for _, feature := range location.InfraTypes {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		contribution := feature.Weight * math.Log1p(feature.Influence)
		if feature.IsCompetitor {
			competition += contribution
		} else {
			positive += contribution
		}
	}

	raw := positive - competition
	return &domain.CalculatedRating{Value: roundRating(logisticRating(raw))}, nil
}

func roundRating(rating float64) float64 {
	rounded := math.Round(rating*10) / 10
	return min(9.9, max(0.1, rounded))
}

func logisticRating(raw float64) float64 {
	if raw >= 0 {
		return 10 / (1 + math.Exp(-raw))
	}
	exp := math.Exp(raw)
	return 10 * exp / (1 + exp)
}
