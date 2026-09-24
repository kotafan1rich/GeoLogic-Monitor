package domain

import (
	"uuid"
)

type LocationRating struct {
	ID                uuid.UUID
	TrackedLocationID uuid.UUID
	CalculatedRating
}

func NewLocationRating(
	trackedLocationID uuid.UUID,
	rating CalculatedRating,
) (*LocationRating, error) {
	return &LocationRating{
		TrackedLocationID: trackedLocationID,
		CalculatedRating:  rating,
	}, nil
}

type LocationRatingHistory struct {
	TrackedLocationID uuid.UUID
	History           []*CalculatedRating
}

func NewLoctionRatingHistory(trackedLocationID uuid.UUID, history *[]*CalculatedRating) *LocationRatingHistory {
	return &LocationRatingHistory{
		TrackedLocationID: trackedLocationID,
		History:           *history,
	}
}
