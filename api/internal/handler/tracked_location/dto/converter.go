package dto

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"

func ToResponse(location domain.TrackedLocation) TrackedLocationResponse {
	return TrackedLocationResponse{
		ID:             location.ID,
		Name:           location.Name,
		BusinessTypeID: location.BusinessTypeID,
		Address:        location.Address,
		Lat:            location.GeoPoint.Lat,
		Lon:            location.GeoPoint.Lng,
	}
}

func ToResponseList(locations []domain.TrackedLocation) []UserTrackedLocationResponse {
	result := make([]UserTrackedLocationResponse, 0, len(locations))
	for _, location := range locations {
		response := UserTrackedLocationResponse{
			TrackedLocationResponse: ToResponse(location),
		}
		if location.LatestRating != nil {
			response.Rating = &location.LatestRating.Value
			response.RatingCalculatedAt = &location.LatestRating.CalculatedAt
		}
		result = append(result, response)
	}
	return result
}

func ToCreatedResponse(location domain.TrackedLocationRating) *CreatedTrackedLocationResponse {
	return &CreatedTrackedLocationResponse{
		TrackedLocationResponse: ToResponse(location.TrackedLocation),
		Rating:                  location.Value,
		RatingCalculatedAt:      location.CalculatedAt,
	}
}

func ToRatingHistoryResponse(history *domain.LocationRatingHistory) []RatingHistoryEntryResponse {
	if history == nil {
		return []RatingHistoryEntryResponse{}
	}

	result := make([]RatingHistoryEntryResponse, 0, len(history.History))
	for _, rating := range history.History {
		result = append(result, RatingHistoryEntryResponse{
			Value:        rating.Value,
			CalculatedAt: rating.CalculatedAt,
		})
	}
	return result
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
