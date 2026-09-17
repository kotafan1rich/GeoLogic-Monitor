package model

import (
	"uuid"

	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

type BusinessType struct {
	basemodel.Base
	InfraTypeID uuid.UUID
	InfraType   InfraType
}

type InfraType struct {
	ID        uuid.UUID
	Slug      string
	Name      string
	Weight    float64
	MaxRadius uint16
}
