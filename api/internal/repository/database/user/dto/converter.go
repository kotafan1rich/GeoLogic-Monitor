package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/user/model"
)

func ToModel(user domain.User) *model.User {
	return &model.User{
		ID:        user.ID,
		MaxUserID: user.MaxUserID,
		MaxChatID: user.MaxChatID,
	}
}

func ToDomain(user model.User) *domain.User {
	return &domain.User{
		ID:        user.ID,
		MaxUserID: user.MaxUserID,
		MaxChatID: user.MaxChatID,
	}
}
