package calculate

import (
	"math"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

type groupedLocationFeatures struct {
	infraType domain.InfraType
	distances []float64
}

func BuildLocationFeatures(
	businessType *domain.BusinessType,
	objects []*domain.InfraObjectDistance,
) *domain.LocationFeatures {
	grouped := make([]groupedLocationFeatures, 0)
	byTypeID := make(map[uuid.UUID]int)

	for _, object := range objects {
		if object == nil || object.Object == nil {
			continue
		}

		typeID := object.Object.Type.ID
		index, exists := byTypeID[typeID]
		if !exists {
			grouped = append(grouped, groupedLocationFeatures{infraType: object.Object.Type})
			index = len(grouped) - 1
			byTypeID[typeID] = index
		}
		grouped[index].distances = append(grouped[index].distances, object.DistanceMeters)
	}

	features := &domain.LocationFeatures{
		BusinessTypeID: businessType.ID,
		InfraTypes:     make([]domain.InfraTypeFeatures, 0, len(grouped)),
	}
	for _, group := range grouped {
		features.InfraTypes = append(
			features.InfraTypes,
			toInfraTypeFeatures(group, businessType.InfraTypeID),
		)
	}

	return features
}

func toInfraTypeFeatures(
	group groupedLocationFeatures,
	competitorTypeID uuid.UUID,
) domain.InfraTypeFeatures {
	nearest := group.distances[0]
	var distanceSum, influence float64
	var count100m, count300m, count500m uint64

	for _, distance := range group.distances {
		nearest = min(nearest, distance)
		distanceSum += distance
		influence += math.Pow(
			math.Max(0, 1-distance/float64(group.infraType.MaxRadius)),
			2,
		)
		if distance <= 100 {
			count100m++
		}
		if distance <= 300 {
			count300m++
		}
		if distance <= 500 {
			count500m++
		}
	}

	mean := distanceSum / float64(len(group.distances))
	return domain.InfraTypeFeatures{
		TypeID:          group.infraType.ID,
		Slug:            group.infraType.Slug,
		Weight:          group.infraType.Weight,
		MaxRadiusMeters: group.infraType.MaxRadius,
		ObjectCount:     uint64(len(group.distances)),
		NearestMeters:   &nearest,
		MeanMeters:      &mean,
		DistanceSum:     distanceSum,
		Influence:       influence,
		Count100m:       count100m,
		Count300m:       count300m,
		Count500m:       count500m,
		IsCompetitor:    group.infraType.ID == competitorTypeID,
	}
}
