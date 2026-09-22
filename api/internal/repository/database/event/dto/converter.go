package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/event/model"
	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

func ToDomain(event *model.Event) *domain.Event {
	return &domain.Event{
		ID:         event.ID,
		Provider:   event.Provider,
		ExternalID: event.ExternalID,
		GeoPoint:   domain.GeoPoint(event.Location),
		Date:       event.Date,
		Info:       event.Info,
		NotifiedAt: event.NotifiedAt,
	}
}

func ToModel(event *domain.Event) *model.Event {
	return &model.Event{
		ID:         event.ID,
		Provider:   event.Provider,
		ExternalID: event.ExternalID,
		Location:   basemodel.GeoPoint(event.GeoPoint),
		Date:       event.Date,
		Info:       event.Info,
		NotifiedAt: event.NotifiedAt,
	}
}
