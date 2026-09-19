package rating

import (
	"math"
	"testing"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

func TestBuildLocationFeatures(t *testing.T) {
	t.Parallel()

	businessTypeID := uuid.New()
	competitorTypeID := uuid.New()
	positiveTypeID := uuid.New()
	objects := []*domain.InfraObjectDistance{
		newObjectDistance(competitorTypeID, "competitor", 2, 500, 100),
		newObjectDistance(positiveTypeID, "metro", 1.5, 500, 50),
		newObjectDistance(competitorTypeID, "competitor", 2, 500, 300),
	}
	businessType := &domain.BusinessType{
		ID:          businessTypeID,
		InfraTypeID: competitorTypeID,
	}

	features := BuildLocationFeatures(businessType, objects)
	if features.BusinessTypeID != businessTypeID {
		t.Fatalf("business type ID: got %v, want %v", features.BusinessTypeID, businessTypeID)
	}
	if len(features.InfraTypes) != 2 {
		t.Fatalf("infra types count: got %d, want 2", len(features.InfraTypes))
	}

	competitor := features.InfraTypes[0]
	if competitor.ObjectCount != 2 || *competitor.NearestMeters != 100 || *competitor.MeanMeters != 200 {
		t.Fatalf("unexpected competitor distances: %+v", competitor)
	}
	if competitor.DistanceSum != 400 || competitor.Count100m != 1 || competitor.Count300m != 2 || competitor.Count500m != 2 {
		t.Fatalf("unexpected competitor aggregates: %+v", competitor)
	}
	if !competitor.IsCompetitor {
		t.Fatal("competitor infra type is not marked as competitor")
	}
	wantInfluence := math.Pow(1-100.0/500, 2) + math.Pow(1-300.0/500, 2)
	if math.Abs(competitor.Influence-wantInfluence) > 1e-9 {
		t.Fatalf("influence: got %v, want %v", competitor.Influence, wantInfluence)
	}

	positive := features.InfraTypes[1]
	if positive.TypeID != positiveTypeID || positive.IsCompetitor {
		t.Fatalf("unexpected positive infra type: %+v", positive)
	}
}

func TestBuildLocationFeaturesWithoutObjects(t *testing.T) {
	t.Parallel()

	features := BuildLocationFeatures(&domain.BusinessType{
		ID:          uuid.New(),
		InfraTypeID: uuid.New(),
	}, nil)
	if features.InfraTypes == nil {
		t.Fatal("infra types must be an empty slice")
	}
	if len(features.InfraTypes) != 0 {
		t.Fatalf("infra types count: got %d, want 0", len(features.InfraTypes))
	}
}

func newObjectDistance(
	typeID uuid.UUID,
	slug string,
	weight float64,
	maxRadius uint16,
	distance float64,
) *domain.InfraObjectDistance {
	return &domain.InfraObjectDistance{
		Object: &domain.InfraObject{
			Type: domain.InfraType{
				ID:        typeID,
				Slug:      slug,
				Weight:    weight,
				MaxRadius: maxRadius,
			},
		},
		DistanceMeters: distance,
	}
}
