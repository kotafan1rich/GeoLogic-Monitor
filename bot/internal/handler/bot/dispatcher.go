package bot

import (
	"context"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type StartHandler interface {
	Start(ctx context.Context, chatID, userID int64)
}

type dispatcher struct {
	startHandler StartHandler
}

func NewDispatcher(startHandler StartHandler) *dispatcher {
	return &dispatcher{
		startHandler: startHandler,
	}
}

func (d *dispatcher) Dispatch(ctx context.Context, update model.Update) {
	switch update.UpdateType {
	case model.UpdateBotStarted:
		d.startHandler.Start(ctx, update.ChatID, update.UserID)

	case model.UpdateMessageCreated:
		if update.GetCommand().Command == "/start" {
			d.startHandler.Start(ctx, update.ChatID, update.UserID)
		}
	}
}
