package domain

import "uuid"

type User struct {
	ID        uuid.UUID
	MaxUserId int64
	MaxChatId int64
}

func NewUser(maxUserId, maxChatId int64) *User {
	return &User{
		MaxUserId: maxUserId,
		MaxChatId: maxChatId,
	}
}
