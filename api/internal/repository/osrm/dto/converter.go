package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/osrm"
)

func ToCoordinate(point *domain.GeoPoint) *osrm.Coordinate {
	return &osrm.Coordinate{
		Lat: point.Lat,
		Lon: point.Lon,
	}
}

func ToCoordinateSlice(points []*domain.GeoPoint) []*osrm.Coordinate {
	result := make([]*osrm.Coordinate, len(points))
	for i, point := range points {
		result[i] = ToCoordinate(point)
	}
	return result
}
