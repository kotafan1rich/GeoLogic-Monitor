package model

import basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"

type InfraType struct {
	basemodel.Base
	Slug      string
	Name      string
	Weight    float64
	MaxRadius uint16
}
