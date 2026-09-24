package repository

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
)

type ApiRepository interface {
	PutUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type Messenger interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}
