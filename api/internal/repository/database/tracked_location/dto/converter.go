package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/tracked_location/model"
)

func ToDomain(trackedLocation model.TrackedLocation) *domain.TrackedLocation {
	result := &domain.TrackedLocation{
		ID:             trackedLocation.ID,
		UserID:         trackedLocation.UserID,
		BusinessTypeID: trackedLocation.BusinessTypeID,
		Name:           trackedLocation.Name,
		Address:        trackedLocation.Address,
		GeoPoint:       domain.GeoPoint(trackedLocation.Location),
		User: domain.User{
			ID:        trackedLocation.User.ID,
			MaxUserID: trackedLocation.User.MaxUserID,
			MaxChatID: trackedLocation.User.MaxChatID,
		},
	}
	if trackedLocation.LatestRating != nil && trackedLocation.RatingCalculatedAt != nil {
		result.LatestRating = &domain.CalculatedRating{
			Value:        *trackedLocation.LatestRating,
			CalculatedAt: *trackedLocation.RatingCalculatedAt,
		}
	}
	return result
}

func ToModel(trackedLocation domain.TrackedLocation) *model.TrackedLocation {
	return &model.TrackedLocation{
		ID:             trackedLocation.ID,
		UserID:         trackedLocation.UserID,
		BusinessTypeID: trackedLocation.BusinessTypeID,
		Name:           trackedLocation.Name,
		Address:        trackedLocation.Address,
		Location: basemodel.GeoPoint{
			Lat: trackedLocation.GeoPoint.Lat,
			Lng: trackedLocation.GeoPoint.Lng,
		},
		User: model.User{
			ID:        trackedLocation.User.ID,
			MaxUserID: trackedLocation.User.MaxUserID,
			MaxChatID: trackedLocation.User.MaxChatID,
		},
	}
}
