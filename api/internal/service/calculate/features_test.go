package calculate

import (
	"testing"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

func TestBuildLocationFeaturesPreservesWalkingDistance(t *testing.T) {
	business := &domain.BusinessType{ID: uuid.New(), InfraType: domain.InfraType{Slug: "restaurant"}}
	features := BuildLocationFeatures(business, []*domain.InfraObjectDistance{
		nil,
		{Object: nil},
		{Object: &domain.InfraObject{Type: domain.InfraType{Slug: "subway", Weight: 999, MaxRadius: 1}}, DistanceMeters: 300},
	})
	if features.BusinessTypeID != business.ID || features.BusinessSlug != "restaurant" || len(features.Objects) != 1 || features.Objects[0].Type != "subway" || features.Objects[0].DistanceMeters != 300 {
		t.Fatalf("unexpected features: %+v", features)
	}
	if len(BuildLocationFeatures(business, nil).Objects) != 0 {
		t.Fatal("empty objects")
	}
}
