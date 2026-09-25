package handler

import (
	"context"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/broker"
)

type NotificationHandler interface {
	Handle(ctx context.Context, message broker.Message) error
}

type WebhookHandler interface {
	Handle(ctx context.Context, update model.Update)
}
