package config

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

type RatingFormula struct {
	Beta                float64                       `yaml:"beta"`
	Gamma               float64                       `yaml:"gamma"`
	ScoreCalibration    float64                       `yaml:"score_calibration"`
	Scale               float64                       `yaml:"scale"`
	Min                 float64                       `yaml:"min"`
	Max                 float64                       `yaml:"max"`
	RoundFactor         float64                       `yaml:"round_factor"`
	BreakdownLimit      int                           `yaml:"breakdown_limit"`
	StreetFactorDefault float64                       `yaml:"street_factor_default"`
	StreetFactorMin     float64                       `yaml:"street_factor_min"`
	StreetFactorMax     float64                       `yaml:"street_factor_max"`
	AvailableTypes      []string                      `yaml:"available_types"`
	InfraTypes          []RatingInfraType             `yaml:"infra_types"`
	Profiles            map[string]map[string]float64 `yaml:"profiles"`
}

type RatingInfraType struct {
	Slug       string  `yaml:"slug"`
	Weight     float64 `yaml:"w"`
	Radius     float64 `yaml:"R"`
	Saturation float64 `yaml:"k"`
}

func LoadRatingFormula(path string) (RatingFormula, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return RatingFormula{}, fmt.Errorf("read rating config: %w", err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var cfg RatingFormula
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("decode rating config: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return cfg, fmt.Errorf("rating config must contain exactly one YAML document")
	}
	return cfg, cfg.Validate()
}

func (c RatingFormula) MaxRadius() float64 {
	var radius float64
	for _, t := range c.InfraTypes {
		radius = max(radius, t.Radius)
	}
	return radius
}

func (c RatingFormula) Validate() error {
	if !finite(c.Beta) || c.Beta < 0 || c.Beta > 1 || !positive(c.Gamma) {
		return fmt.Errorf("rating: beta must be in [0,1] and gamma must be positive and finite")
	}
	if !positive(c.ScoreCalibration) {
		return fmt.Errorf("rating: score calibration must be positive and finite")
	}
	// Output must fit the existing NUMERIC(2,1) history column.
	if !positive(c.Scale) || !positive(c.Min) || !finite(c.Max) || c.Min < 0.1 || c.Max > 9.9 || c.Min > c.Max || c.Max > c.Scale || c.RoundFactor != 10 || c.BreakdownLimit <= 0 {
		return fmt.Errorf("rating: invalid output scale, bounds, rounding or breakdown limit")
	}
	if math.Round(c.Min*c.RoundFactor)/c.RoundFactor != c.Min || math.Round(c.Max*c.RoundFactor)/c.RoundFactor != c.Max {
		return fmt.Errorf("rating: bounds must be representable to one decimal place")
	}
	if !finite(c.StreetFactorMin) || !finite(c.StreetFactorMax) || !finite(c.StreetFactorDefault) || c.StreetFactorMin < 0.85 || c.StreetFactorMax > 1 || c.StreetFactorMin > c.StreetFactorMax || c.StreetFactorDefault < c.StreetFactorMin || c.StreetFactorDefault > c.StreetFactorMax {
		return fmt.Errorf("rating: invalid street factor bounds or default")
	}
	if len(c.InfraTypes) == 0 {
		return fmt.Errorf("rating: infra_types must not be empty")
	}
	known := make(map[string]bool, len(c.InfraTypes))
	for _, t := range c.InfraTypes {
		if strings.TrimSpace(t.Slug) != t.Slug || t.Slug == "" || known[t.Slug] || !finite(t.Weight) || t.Weight < 0 || !positive(t.Radius) || !positive(t.Saturation) {
			return fmt.Errorf("rating: invalid or duplicate infra type %q", t.Slug)
		}
		known[t.Slug] = true
	}
	available := make(map[string]bool, len(c.AvailableTypes))
	for _, slug := range c.AvailableTypes {
		if !known[slug] || available[slug] {
			return fmt.Errorf("rating: unknown or duplicate available type %q", slug)
		}
		available[slug] = true
	}
	for business, profile := range c.Profiles {
		if !known[business] {
			return fmt.Errorf("rating: unknown business profile %q", business)
		}
		for slug, rel := range profile {
			if !known[slug] || !finite(rel) || rel < 0 {
				return fmt.Errorf("rating: invalid relevance for %q/%q", business, slug)
			}
		}
	}
	return nil
}

func finite(v float64) bool   { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func positive(v float64) bool { return finite(v) && v > 0 }
