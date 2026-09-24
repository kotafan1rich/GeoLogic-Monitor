package dto

import "uuid"

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	MaxUserID int64     `json:"max_user_id"`
	MaxChatID int64     `json:"max_chat_id"`
}
