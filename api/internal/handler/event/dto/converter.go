package dto

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"

func ToResponse(event *domain.Event) *EventResponse {
	return &EventResponse{
		ID:         event.ID,
		Provider:   event.Provider,
		ExternalID: event.ExternalID,
		Lat:        event.GeoPoint.Lat,
		Lon:        event.GeoPoint.Lng,
		Date:       event.Date,
		Info:       event.Info,
		NotifiedAt: event.NotifiedAt,
	}
}

func ToResponseList(events []*domain.Event) []*EventResponse {
	result := make([]*EventResponse, 0, len(events))
	for _, event := range events {
		result = append(result, ToResponse(event))
	}
	return result
}
