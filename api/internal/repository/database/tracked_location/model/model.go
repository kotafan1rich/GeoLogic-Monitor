package model

import (
	"time"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

type TrackedLocation struct {
	model.Base
	UserID             uuid.UUID
	BusinessTypeID     uuid.UUID
	Address            string
	Location           model.GeoPoint
	User               User
	LatestRating       *float64
	RatingCalculatedAt *time.Time
}

type User struct {
	model.Base
	MaxUserID int64
	MaxChatID int64
}
