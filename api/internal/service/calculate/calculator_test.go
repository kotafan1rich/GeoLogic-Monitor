package calculate

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

func testConfig(t testing.TB) config.RatingFormula {
	t.Helper()
	cfg, err := config.LoadRatingFormula("../../../rating.yml")
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func assess(t testing.TB, cfg config.RatingFormula, in domain.LocationFeatures) *domain.RatingAssessment {
	t.Helper()
	out, err := NewFormulaCalculator(cfg).Calculate(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestRatingProperties(t *testing.T) {
	cfg := testConfig(t)
	input := domain.LocationFeatures{BusinessSlug: "restaurant"}
	empty := assess(t, cfg, input)
	if empty.Rating != 0.1 || empty.Confidence != 1 {
		t.Fatalf("empty: %+v", empty)
	}
	metro := func(d float64) *domain.RatingAssessment {
		return assess(t, cfg, domain.LocationFeatures{BusinessSlug: "restaurant", Objects: []domain.RatingObjectDistance{{Type: "subway", DistanceMeters: d}}})
	}
	if metro(300).Rating >= 1.5 || metro(100).Rating <= metro(600).Rating {
		t.Fatal("metro saturation or distance monotonicity failed")
	}
	if metro(800).Rating != empty.Rating || metro(900).Rating != empty.Rating {
		t.Fatal("radius boundary failed")
	}
	input.Objects = []domain.RatingObjectDistance{{Type: "subway", DistanceMeters: 100}}
	previous := assess(t, cfg, input).Rating
	for range 50 {
		input.Objects = append(input.Objects, domain.RatingObjectDistance{Type: "restaurant", DistanceMeters: 50})
		next := assess(t, cfg, input)
		if next.Rating > previous || next.CompetitionPenalty < 0 || next.CompetitionPenalty > cfg.Beta {
			t.Fatal("competitor raised rating or invalid penalty")
		}
		previous = next.Rating
	}
	for _, typ := range cfg.InfraTypes {
		if typ.Slug == "restaurant" {
			continue
		}
		previous = assess(t, cfg, input).Rating
		input.Objects = append(input.Objects, domain.RatingObjectDistance{Type: typ.Slug, DistanceMeters: 0})
		if assess(t, cfg, input).Rating < previous {
			t.Fatalf("%s lowered rating", typ.Slug)
		}
	}
	cafes := func(n int) float64 {
		in := domain.LocationFeatures{BusinessSlug: "restaurant"}
		for range n {
			in.Objects = append(in.Objects, domain.RatingObjectDistance{Type: "cafe", DistanceMeters: 100})
		}
		return assess(t, cfg, in).Rating
	}
	if cafes(50)-cafes(10) >= cafes(10)-cafes(0) {
		t.Fatal("cafe saturation failed")
	}
	if !reflect.DeepEqual(assess(t, cfg, input), assess(t, cfg, input)) {
		t.Fatal("non-deterministic output")
	}
}

func TestAvailabilityProfilesAndBreakdown(t *testing.T) {
	cfg := testConfig(t)
	cfg.AvailableTypes = []string{"subway", "cafe"}
	input := domain.LocationFeatures{BusinessSlug: "restaurant", Objects: []domain.RatingObjectDistance{{Type: "subway", DistanceMeters: 100}}}
	initial := assess(t, cfg, input)
	input.Objects = append(input.Objects, domain.RatingObjectDistance{Type: "hotel", DistanceMeters: 0}, domain.RatingObjectDistance{Type: "restaurant", DistanceMeters: 0})
	if !reflect.DeepEqual(initial, assess(t, cfg, input)) {
		t.Fatal("unavailable type changed assessment")
	}
	cfg.Profiles = map[string]map[string]float64{"restaurant": {"cafe": 0}}
	out := assess(t, cfg, input)
	if out.Rating <= initial.Rating {
		t.Fatal("configured profile not applied")
	}
	input.Profile = map[string]float64{"subway": 0}
	out = assess(t, cfg, input)
	if out.Rating != cfg.Min || out.Confidence != 0 || len(out.Breakdown) != 0 {
		t.Fatalf("zero denominator: %+v", out)
	}
	input.Profile = nil
	street := 0.85
	input.StreetFactor = &street
	if assess(t, cfg, input).Rating >= assess(t, cfg, domain.LocationFeatures{BusinessSlug: "restaurant", Objects: input.Objects}).Rating {
		t.Fatal("street factor not applied")
	}
	cfg.AvailableTypes = nil
	out = assess(t, cfg, input)
	if out.Rating != cfg.Min || out.Confidence != 0 {
		t.Fatal("empty availability")
	}

	cfg = testConfig(t)
	cfg.InfraTypes = []config.RatingInfraType{
		{Slug: "restaurant", Weight: 3, Radius: 300, Saturation: 1},
		{Slug: "b", Weight: 1, Radius: 300, Saturation: 1},
		{Slug: "a", Weight: 1, Radius: 300, Saturation: 1},
		{Slug: "c", Weight: 2, Radius: 300, Saturation: 1},
		{Slug: "zero", Weight: 0, Radius: 300, Saturation: 1},
	}
	cfg.AvailableTypes = []string{"restaurant", "a", "b", "zero"}
	input = domain.LocationFeatures{BusinessSlug: "restaurant", Objects: []domain.RatingObjectDistance{{Type: "b"}, {Type: "a"}, {Type: "zero"}}}
	out = assess(t, cfg, input)
	if out.Confidence != 0.5 || len(out.Breakdown) != 2 || out.Breakdown[0].Type != "a" || out.Breakdown[1].Type != "b" {
		t.Fatalf("confidence/order: %+v", out)
	}
	if math.Abs(out.Breakdown[0].Contribution-(-math.Expm1(-1)/2)) > 1e-12 {
		t.Fatal("contribution is not normalized")
	}
	cfg.BreakdownLimit = 1
	if len(assess(t, cfg, input).Breakdown) != 1 {
		t.Fatal("breakdown limit")
	}
	cfg.AvailableTypes = append(cfg.AvailableTypes, "c")
	if assess(t, cfg, input).Confidence != 1 || assess(t, cfg, input).Rating >= out.Rating {
		t.Fatal("availability must affect denominator and confidence")
	}
}

func TestBoundsAndInvalidInputs(t *testing.T) {
	cfg := testConfig(t)
	in := domain.LocationFeatures{BusinessSlug: "restaurant"}
	for _, typ := range cfg.InfraTypes {
		if typ.Slug == "restaurant" {
			continue
		}
		for range 100 {
			in.Objects = append(in.Objects, domain.RatingObjectDistance{Type: typ.Slug})
		}
	}
	result := assess(t, cfg, in)
	if result.Rating != 9.9 || len(result.Breakdown) != 5 {
		t.Fatalf("upper clamp: %+v", result)
	}
	for _, d := range []float64{-1, math.NaN(), math.Inf(1)} {
		_, err := NewFormulaCalculator(cfg).Calculate(context.Background(), domain.LocationFeatures{BusinessSlug: "restaurant", Objects: []domain.RatingObjectDistance{{Type: "subway", DistanceMeters: d}}})
		if err == nil {
			t.Fatalf("accepted distance %v", d)
		}
	}
	for _, street := range []float64{0, 0.84, 1.01, math.NaN()} {
		_, err := NewFormulaCalculator(cfg).Calculate(context.Background(), domain.LocationFeatures{BusinessSlug: "restaurant", StreetFactor: &street})
		if err == nil {
			t.Fatalf("accepted street factor %v", street)
		}
	}
	for _, profile := range []map[string]float64{{"subway": -1}, {"subway": math.NaN()}, {"unknown": 1}} {
		_, err := NewFormulaCalculator(cfg).Calculate(context.Background(), domain.LocationFeatures{BusinessSlug: "restaurant", Profile: profile})
		if err == nil {
			t.Fatal("accepted invalid profile")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := NewFormulaCalculator(cfg).Calculate(ctx, in); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
}

func BenchmarkFormulaCalculator1000Objects(b *testing.B) {
	cfg := testConfig(b)
	in := domain.LocationFeatures{BusinessSlug: "restaurant", Objects: make([]domain.RatingObjectDistance, 1000)}
	for i := range in.Objects {
		typ := cfg.InfraTypes[i%len(cfg.InfraTypes)]
		in.Objects[i] = domain.RatingObjectDistance{Type: typ.Slug, DistanceMeters: float64(i % 200)}
	}
	calculator := NewFormulaCalculator(cfg)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := calculator.Calculate(context.Background(), in); err != nil {
			b.Fatal(err)
		}
	}
}
