package domain

import "uuid"

type TrackedLocation struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	BusinessTypeID uuid.UUID
	Address        string
	GeoPoint       GeoPoint
	User           User
	LatestRating   *CalculatedRating
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

type TrackedLocationRating struct {
	TrackedLocation
	CalculatedRating
}

func NewTrackedLocationRating(location *TrackedLocation, rating *CalculatedRating) *TrackedLocationRating {
	return &TrackedLocationRating{
		TrackedLocation:  *location,
		CalculatedRating: *rating,
	}
}

type MonitoringLocation struct {
	TrackedLocation
	MaxChatID int64
}

func NewMonitoringLocation(location *TrackedLocation, maxChatID int64) *MonitoringLocation {
	return &MonitoringLocation{
		TrackedLocation: *location,
		MaxChatID:       maxChatID,
	}
}
