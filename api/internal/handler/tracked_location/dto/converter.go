package dto

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"

func ToResponse(location domain.TrackedLocation) TrackedLocationResponse {
	return TrackedLocationResponse{
		ID:             location.ID,
		BusinessTypeID: location.BusinessTypeID,
		Address:        location.Address,
		Lat:            location.GeoPoint.Lat,
		Lon:            location.GeoPoint.Lng,
	}
}

func ToResponseList(locations []domain.TrackedLocation) []TrackedLocationResponse {
	result := make([]TrackedLocationResponse, 0, len(locations))
	for _, location := range locations {
		result = append(result, ToResponse(location))
	}
	return result
}

func ToCreatedResponse(location domain.TrackedLocationRating) *CreatedTrackedLocationResponse {
	return &CreatedTrackedLocationResponse{
		TrackedLocationResponse: ToResponse(location.TrackedLocation),
		Rating:                  location.Value,
	}
}

func ToMonitoringResponseList(locations []domain.MonitoringLocation) []MonitoringLocationResponse {
	result := make([]MonitoringLocationResponse, 0, len(locations))
	for _, location := range locations {
		result = append(result, MonitoringLocationResponse{
			TrackedLocationResponse: ToResponse(location.TrackedLocation),
			MaxChatID:               location.MaxChatID,
		})
	}
	return result
}
