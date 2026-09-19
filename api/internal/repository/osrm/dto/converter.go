package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/osrm"
)

func ToCoordinate(point *domain.GeoPoint) *osrm.Coordinate {
	return &osrm.Coordinate{
		Lat: point.Lat,
		Lon: point.Lng,
	}
}

func ToCoordinateSlice(infras []*domain.InfraObject) []*osrm.Coordinate {
	result := make([]*osrm.Coordinate, len(infras))
	for i, infra := range infras {
		result[i] = ToCoordinate(&infra.GeoPoint)
	}
	return result
}
