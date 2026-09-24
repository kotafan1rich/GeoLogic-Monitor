package dto

import (
	"errors"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

func ToResponse(distances []*float64) *WalkingDistancesResponse {
	return &WalkingDistancesResponse{Distances: distances}
}

func ToGeoPoint(point GeoPointRequest) (*domain.GeoPoint, error) {
	if point.Lat == nil || point.Lon == nil {
		return nil, errors.New("lat and lon are required")
	}
	return domain.NewGeoPoint(*point.Lat, *point.Lon)
}
