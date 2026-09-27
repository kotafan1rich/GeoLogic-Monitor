package calculate

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"

// BuildLocationFeatures preserves OSRM distances. Configured radii are applied by the calculator.
func BuildLocationFeatures(
	businessType *domain.BusinessType,
	objects []*domain.InfraObjectDistance,
) *domain.LocationFeatures {
	features := &domain.LocationFeatures{
		BusinessTypeID: businessType.ID,
		BusinessSlug:   businessType.InfraType.Slug,
		Objects:        make([]domain.RatingObjectDistance, 0, len(objects)),
	}
	for _, object := range objects {
		if object == nil || object.Object == nil {
			continue
		}
		features.Objects = append(features.Objects, domain.RatingObjectDistance{
			Type:           object.Object.Type.Slug,
			DistanceMeters: object.DistanceMeters,
		})
	}
	return features
}
