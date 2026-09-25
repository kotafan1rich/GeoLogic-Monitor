package service

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain/notification"
)

type OnboardingService interface {
	RegisterUser(ctx context.Context, maxUserID, maxChatID int64) error
}

type NotificationService interface {
	Notify(ctx context.Context, notice *notification.Notification) error
}
