package domain

import (
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
)

type InfraObjectDistance struct {
	Object         *InfraObject
	DistanceMeters float64
}

type InfraObject struct {
	ID         uuid.UUID
	ExternalID string
	TypeID     uuid.UUID
	GeoPoint   GeoPoint
	Address    string
	Name       *string
	Type       InfraType
}

func NewInfraObject(
	externalID string,
	typeID uuid.UUID,
	geoPoint *GeoPoint,
	address string,
	name *string,
) (*InfraObject, error) {
	if externalID == "" {
		return nil, errs.ErrInvalidExternalID
	}
	if address == "" {
		return nil, errs.ErrInvalidAddress
	}

	return &InfraObject{
		ExternalID: externalID,
		TypeID:     typeID,
		GeoPoint:   *geoPoint,
		Address:    address,
		Name:       name,
	}, nil
}
