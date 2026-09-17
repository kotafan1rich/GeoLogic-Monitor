package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/tracked_location/model"
)

func ToDomain(trackedLocation model.TrackedLocation) *domain.TrackedLocation {
	return &domain.TrackedLocation{
		ID:             trackedLocation.ID,
		UserID:         trackedLocation.UserID,
		BusinessTypeID: trackedLocation.BusinessTypeID,
		Address:        trackedLocation.Address,
		GeoPoint:       domain.GeoPoint(trackedLocation.Location),
	}
}

func ToModel(trackedLocation domain.TrackedLocation) *model.TrackedLocation {
	return &model.TrackedLocation{
		ID:             trackedLocation.ID,
		UserID:         trackedLocation.UserID,
		BusinessTypeID: trackedLocation.BusinessTypeID,
		Address:        trackedLocation.Address,
		Location: basemodel.GeoPoint{
			Lat: trackedLocation.GeoPoint.Lat,
			Lng: trackedLocation.GeoPoint.Lng,
		},
	}
}
