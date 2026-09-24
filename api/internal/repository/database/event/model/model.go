package model

import (
	"time"

	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

type Event struct {
	basemodel.Base
	Provider   string
	ExternalID string
	Location   basemodel.GeoPoint
	Date       time.Time
	Info       *string
	NotifiedAt *time.Time
}
