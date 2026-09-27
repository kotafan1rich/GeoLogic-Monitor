package trackedlocation

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/calculate"
	ratingservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/rating"
)

func TestRatingUsesOSRMDistancesAndConfiguredRadius(t *testing.T) {
	for _, tc := range []struct {
		name     string
		distance *float64
		failure  error
		missing  bool
	}{
		{"walking distance", floatPointer(300), nil, false},
		{"beyond YAML radius", floatPointer(801), nil, false},
		{"no route", nil, nil, false},
		{"provider unavailable", nil, errors.New("OSRM unavailable"), false},
		{"invalid distance count", nil, nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(&fakeTrackedLocationRepository{})
			cfg := ratingConfig(t)
			s.ratingRadius = cfg.MaxRadius()
			infra := &recordingInfra{objects: []*domain.InfraObject{{
				GeoPoint: domain.GeoPoint{Lat: 59.93, Lon: 30.32},
				Type:     domain.InfraType{Slug: "subway", Weight: 1000, MaxRadius: 1},
			}}}
			s.infraService = infra
			osrm := &recordingOSRM{distances: []*float64{tc.distance}, err: tc.failure}
			if tc.missing {
				osrm.distances = nil
			}
			s.osrmService = osrm
			history := &ratingHistoryRepository{}
			s.ratingService = ratingservice.NewService(testLogger(), calculate.NewFormulaCalculator(cfg), history)
			tx := &checkedTransaction{}
			s.txManager = tx
			result, err := s.Create(context.Background(), uuid.New(), uuid.New(), "name", "address", 59.93, 30.32)
			if infra.calls != 1 || infra.radius != 800 || osrm.calls != 1 || len(osrm.destinations) != 1 {
				t.Fatal("wrong query count/radius or OSRM inputs")
			}
			if tc.failure != nil || tc.missing {
				if err == nil || len(history.rows) != 0 || tx.committed {
					t.Fatal("failed OSRM must prevent saving and committing")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			// Coordinates coincide, but OSRM says 300m; DB radius=1 must not filter it.
			want, err := calculate.NewFormulaCalculator(cfg).Calculate(context.Background(), domain.LocationFeatures{
				BusinessSlug: "restaurant", Objects: func() []domain.RatingObjectDistance {
					if tc.distance == nil {
						return nil
					}
					return []domain.RatingObjectDistance{{Type: "subway", DistanceMeters: *tc.distance}}
				}(),
			})
			if err != nil {
				t.Fatal(err)
			}
			if result.Value != want.Rating || len(history.rows) != 1 || result.CalculatedAt.IsZero() || !tx.committed {
				t.Fatalf("unexpected rating/history: %+v", result)
			}
			if tc.distance != nil && *tc.distance == 300 && result.Value <= 0.1 {
				t.Fatal("DB radius incorrectly excluded subway")
			}
		})
	}
}

func floatPointer(v float64) *float64 { return &v }

type recordingInfra struct {
	objects []*domain.InfraObject
	calls   int
	radius  float64
}

func (r *recordingInfra) Near(_ context.Context, _ *domain.GeoPoint, radius float64) ([]*domain.InfraObject, error) {
	r.calls++
	r.radius = radius
	return r.objects, nil
}

type recordingOSRM struct {
	distances    []*float64
	err          error
	calls        int
	destinations []*domain.GeoPoint
}

func (r *recordingOSRM) WalkingDistances(_ context.Context, _ *domain.GeoPoint, dst []*domain.GeoPoint) ([]*float64, error) {
	r.calls++
	r.destinations = dst
	return r.distances, r.err
}
