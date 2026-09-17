package domain

import "uuid"

type TrackedLocation struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	GeoPoint GeoPoint
}

func NewTrackedLocation(userID uuid.UUID, geopoint *GeoPoint) *TrackedLocation {
	return &TrackedLocation{
		UserID:   userID,
		GeoPoint: *geopoint,
	}
}
