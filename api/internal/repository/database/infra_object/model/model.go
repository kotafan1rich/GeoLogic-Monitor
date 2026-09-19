package model

import (
	"uuid"

	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

type InfraObject struct {
	basemodel.Base
	ExternalID string
	TypeID     uuid.UUID
	Location   basemodel.GeoPoint
	Address    string
	Name       *string
	Type       InfraType
}

type InfraType struct {
	ID        uuid.UUID
	Slug      string
	Name      string
	Weight    float64
	MaxRadius uint16
}
