package model

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"

type User struct {
	model.Base
	MaxUserID int64
	MaxChatID int64
}
