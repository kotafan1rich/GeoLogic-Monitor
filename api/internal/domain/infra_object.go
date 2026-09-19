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
	ID       uuid.UUID
	TypeID   uuid.UUID
	GeoPoint GeoPoint
	Address  string
	Name     *string
	Type     InfraType
}

func NewInfraObject(
	id uuid.UUID,
	typeID uuid.UUID,
	geopoint *GeoPoint,
	address string,
	name *string,
) (*InfraObject, error) {
	if address == "" {
		return nil, errs.ErrInvalidAddress
	}

	return &InfraObject{
		ID:       id,
		TypeID:   typeID,
		GeoPoint: *geopoint,
		Address:  address,
		Name:     name,
	}, nil
}
