package domain

import (
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
)

type User struct {
	ID        uuid.UUID
	MaxUserID int64
	MaxChatID int64
}

func NewUser(maxUserID, maxChatID int64) (*User, error) {
	if maxUserID <= 0 {
		return nil, errs.ErrInvalidMaxUserID
	}

	return &User{
		MaxUserID: maxUserID,
		MaxChatID: maxChatID,
	}, nil
}

func (u *User) UpdateMaxChatID(maxChatID int64) {
	u.MaxChatID = maxChatID
}
