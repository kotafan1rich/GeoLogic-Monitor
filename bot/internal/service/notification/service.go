package notification

import (
	"context"
	"fmt"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain/notification"
)

type Messenger interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type service struct {
	messenger Messenger
}

func New(messenger Messenger) *service {
	return &service{
		messenger: messenger,
	}
}

func (s *service) Notify(
	ctx context.Context,
	notice *notification.Notification,
) error {
	message, err := formatMessage(notice)
	if err != nil {
		return err
	}
	err = s.messenger.SendMessage(ctx, notice.MaxChatID, message)
	if err != nil {
		return fmt.Errorf("send notification to MAX: %w", err)
	}

	return nil
}
