package bot

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/logger"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type StartHandler interface {
	Start(ctx context.Context, chatID, userID int64)
}

type dispatcher struct {
	log          *logger.Logger
	startHandler StartHandler
}

func NewDispatcher(log *logger.Logger, startHandler StartHandler) *dispatcher {
	return &dispatcher{
		log:          log,
		startHandler: startHandler,
	}
}

func (d *dispatcher) Dispatch(ctx context.Context, update model.Update) {
	d.log.InfoContext(
		ctx,
		"message received",
		"chat_id", update.ChatID,
	)
	switch update.UpdateType {
	case model.UpdateBotStarted:
		d.startHandler.Start(ctx, update.ChatID, update.UserID)

	case model.UpdateMessageCreated:
		if update.GetCommand().Command == "/start" {
			d.startHandler.Start(ctx, update.ChatID, update.UserID)
		}
	}
}
