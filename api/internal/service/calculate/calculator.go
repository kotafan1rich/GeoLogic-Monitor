package calculate

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

type FormulaCalculator struct{ cfg config.RatingFormula }

// The configuration and input are read-only. Calculate has no clock, network or database access.
func NewFormulaCalculator(cfg config.RatingFormula) FormulaCalculator {
	return FormulaCalculator{cfg: cfg}
}

func (c FormulaCalculator) Calculate(ctx context.Context, location domain.LocationFeatures) (*domain.RatingAssessment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := c.cfg.Validate(); err != nil {
		return nil, err
	}
	types := make(map[string]config.RatingInfraType, len(c.cfg.InfraTypes))
	for _, t := range c.cfg.InfraTypes {
		types[t.Slug] = t
	}
	if _, ok := types[location.BusinessSlug]; !ok {
		return nil, fmt.Errorf("unknown business slug %q", location.BusinessSlug)
	}
	for slug, rel := range location.Profile {
		if _, ok := types[slug]; !ok || !finite(rel) || rel < 0 {
			return nil, fmt.Errorf("invalid profile relevance for %q", slug)
		}
	}
	street := c.cfg.StreetFactorDefault
	if location.StreetFactor != nil {
		street = *location.StreetFactor
	}
	if !finite(street) || street < c.cfg.StreetFactorMin || street > c.cfg.StreetFactorMax {
		return nil, fmt.Errorf("invalid street factor")
	}
	available := make(map[string]bool, len(c.cfg.AvailableTypes))
	for _, slug := range c.cfg.AvailableTypes {
		available[slug] = true
	}
	influence := make(map[string]float64, len(types))
	for _, object := range location.Objects {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		t, ok := types[object.Type]
		if !ok || !available[object.Type] {
			continue
		}
		d := object.DistanceMeters
		if !finite(d) || d < 0 {
			return nil, fmt.Errorf("invalid walking distance for %q", object.Type)
		}
		factor := max(0, 1-d/t.Radius)
		influence[object.Type] += factor * factor
	}
	result := &domain.RatingAssessment{Rating: c.cfg.Min, Breakdown: make([]domain.RatingContribution, 0)}
	var denominator, totalWeight, positive float64
	for _, t := range c.cfg.InfraTypes {
		sat := -math.Expm1(-influence[t.Slug] / t.Saturation)
		if t.Slug == location.BusinessSlug {
			if available[t.Slug] {
				result.CompetitionPenalty = c.cfg.Beta * sat
			}
			continue
		}
		relevance := 1.0
		if rel, ok := c.cfg.Profiles[location.BusinessSlug][t.Slug]; ok {
			relevance = rel
		}
		if rel, ok := location.Profile[t.Slug]; ok {
			relevance = rel
		}
		weight := t.Weight * relevance
		totalWeight += weight
		if !available[t.Slug] {
			continue
		}
		denominator += weight
		contribution := weight * sat
		positive += contribution
		if contribution > 0 {
			result.Breakdown = append(result.Breakdown, domain.RatingContribution{
				Type: t.Slug, Weight: weight, Sat: sat, Contribution: contribution,
			})
		}
	}
	if !finite(totalWeight) || !finite(positive) {
		return nil, fmt.Errorf("rating weights overflow")
	}
	if denominator == 0 {
		return result, nil
	}
	result.Confidence = denominator / totalWeight
	score := positive / denominator
	score = min(1, score+c.cfg.ScoreGain*score*score)
	score *= (1 - result.CompetitionPenalty) * street
	rating := c.cfg.Scale * math.Pow(min(1, max(0, score)), c.cfg.Gamma)
	result.Rating = min(c.cfg.Max, max(c.cfg.Min, math.Round(rating*c.cfg.RoundFactor)/c.cfg.RoundFactor))
	for i := range result.Breakdown {
		result.Breakdown[i].Contribution /= denominator
	}
	sort.Slice(result.Breakdown, func(i, j int) bool {
		a, b := result.Breakdown[i], result.Breakdown[j]
		if a.Contribution == b.Contribution {
			return a.Type < b.Type
		}
		return a.Contribution > b.Contribution
	})
	if len(result.Breakdown) > c.cfg.BreakdownLimit {
		result.Breakdown = result.Breakdown[:c.cfg.BreakdownLimit]
	}
	return result, ctx.Err()
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
