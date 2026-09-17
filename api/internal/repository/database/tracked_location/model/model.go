package model

import (
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

type TrackedLocation struct {
	model.Base
	UserID         uuid.UUID
	BusinessTypeID uuid.UUID
	Address        string
	Location       model.GeoPoint
}
