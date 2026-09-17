package domain

import "uuid"

type TrackedLocation struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	BusinessTypeID uuid.UUID
	Address        string
	GeoPoint       GeoPoint
}

func NewTrackedLocation(
	userID uuid.UUID,
	businessTypeID uuid.UUID,
	address string,
	geopoint *GeoPoint,
) *TrackedLocation {
	return &TrackedLocation{
		UserID:         userID,
		BusinessTypeID: businessTypeID,
		Address:        address,
		GeoPoint:       *geopoint,
	}
}
