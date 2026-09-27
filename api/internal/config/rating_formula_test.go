package config

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/caarlos0/env/v9"
)

func TestRatingConfigValidation(t *testing.T) {
	cfg, err := LoadRatingFormula("../../rating.yml")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.InfraTypes) != 40 || cfg.MaxRadius() != 800 {
		t.Fatal("default config mismatch")
	}
	cases := map[string]func(*RatingFormula){
		"duplicate slug":      func(c *RatingFormula) { c.InfraTypes = append(c.InfraTypes, c.InfraTypes[0]) },
		"unknown available":   func(c *RatingFormula) { c.AvailableTypes = append(c.AvailableTypes, "unknown") },
		"duplicate available": func(c *RatingFormula) { c.AvailableTypes = append(c.AvailableTypes, c.AvailableTypes[0]) },
		"nan weight":          func(c *RatingFormula) { c.InfraTypes[0].Weight = math.NaN() },
		"negative weight":     func(c *RatingFormula) { c.InfraTypes[0].Weight = -1 },
		"zero radius":         func(c *RatingFormula) { c.InfraTypes[0].Radius = 0 },
		"infinite k":          func(c *RatingFormula) { c.InfraTypes[0].Saturation = math.Inf(1) },
		"bad beta":            func(c *RatingFormula) { c.Beta = 1.1 },
		"bad gamma":           func(c *RatingFormula) { c.Gamma = 0 },
		"bad bounds":          func(c *RatingFormula) { c.Max = 10 },
		"fractional bounds":   func(c *RatingFormula) { c.Min = 0.15 },
		"bad street":          func(c *RatingFormula) { c.StreetFactorDefault = 0.8 },
		"unknown profile":     func(c *RatingFormula) { c.Profiles = map[string]map[string]float64{"unknown": {"subway": 1}} },
		"negative relevance":  func(c *RatingFormula) { c.Profiles = map[string]map[string]float64{"restaurant": {"subway": -1}} },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			c, err := LoadRatingFormula("../../rating.yml")
			if err != nil {
				t.Fatal(err)
			}
			mutate(&c)
			if c.Validate() == nil {
				t.Fatal("accepted invalid config")
			}
		})
	}
	cfg.AvailableTypes = []string{}
	if err := cfg.Validate(); err != nil {
		t.Fatal("empty availability must be valid", err)
	}
	var envCfg Rating
	if err := env.Parse(&envCfg); err != nil {
		t.Fatal("formula must be skipped by env parser", err)
	}
}

func TestRatingConfigRejectsUnknownYAML(t *testing.T) {
	data, err := os.ReadFile("../../rating.yml")
	if err != nil {
		t.Fatal(err)
	}
	for _, suffix := range []string{"\nunknown_field: 1\n", "\n---\nbeta: 0.4\n"} {
		path := filepath.Join(t.TempDir(), "rating.yml")
		if err := os.WriteFile(path, append(data, []byte(suffix)...), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadRatingFormula(path); err == nil {
			t.Fatal("accepted malformed configuration")
		}
	}
}
