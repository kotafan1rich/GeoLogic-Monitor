package domain

import (
	"time"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
)

type Event struct {
	ID         uuid.UUID
	Provider   string
	ExternalID string
	GeoPoint   GeoPoint
	Date       time.Time
	Info       *string
	NotifiedAt *time.Time
}

func NewEvent(
	provider string,
	externalID string,
	geopoint *GeoPoint,
	date time.Time,
	info *string,
) (*Event, error) {
	if provider == "" {
		return nil, errs.ErrInvalidEventProvider
	}
	if externalID == "" {
		return nil, errs.ErrInvalidEventExternalID
	}

	return &Event{
		Provider:   provider,
		ExternalID: externalID,
		GeoPoint:   *geopoint,
		Date:       date,
		Info:       info,
	}, nil
}
