package dto

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"

func ToResponse(user domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		MaxUserID: user.MaxUserID,
		MaxChatID: user.MaxChatID,
	}
}
