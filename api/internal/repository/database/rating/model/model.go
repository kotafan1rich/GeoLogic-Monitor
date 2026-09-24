package model

import (
	"time"
	"uuid"

	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

type LocationRating struct {
	basemodel.Base
	TrackedLocationID uuid.UUID
	Value             float64
	CalculatedAt      time.Time
}
