package api

import "uuid"

type UpsertUserRequest struct {
	MaxChatId int64 `json:"max_chat_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	MaxUserId int64     `json:"max_user_id"`
	MaxChatId int64     `json:"max_chat_id"`
}
