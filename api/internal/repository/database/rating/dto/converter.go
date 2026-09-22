package dto

import (
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/rating/model"
)

func ToDomain(locationRating model.LocationRating) *domain.LocationRating {
	return &domain.LocationRating{
		ID:                locationRating.ID,
		TrackedLocationID: locationRating.TrackedLocationID,
		Value:             locationRating.Value,
		CalculatedAt:      locationRating.CalculatedAt,
	}
}

func ToModel(locationRating domain.LocationRating) *model.LocationRating {
	return &model.LocationRating{
		ID:                locationRating.ID,
		TrackedLocationID: locationRating.TrackedLocationID,
		Value:             locationRating.Value,
		CalculatedAt:      locationRating.CalculatedAt,
	}
}

func ToDomainHistory(
	trackedLocationID uuid.UUID,
	history []*model.LocationRating,
) *domain.LocationRatingHistory {
	result := domain.LocationRatingHistory{
		TrackedLocationID: trackedLocationID,
		History:           make([]*domain.CalculatedRating, 0, len(history)),
	}
	for i := range history {
		result.History = append(result.History, &domain.CalculatedRating{
			Value:        history[i].Value,
			CalculatedAt: history[i].CalculatedAt,
		})
	}
	return &result
}
