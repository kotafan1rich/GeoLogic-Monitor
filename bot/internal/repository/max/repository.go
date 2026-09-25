package max

import (
	"context"
	"errors"
	"fmt"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/errs"
	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
)

type repository struct {
	api *maxbot.Api
}

func NewRepository(api *maxbot.Api) *repository {
	return &repository{
		api: api,
	}
}

func (r *repository) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := r.api.Messages.Send(
		ctx,
		maxbot.NewMessage().SetChat(chatID).SetText(text),
	)
	if err != nil {
		var apiErr *maxbot.Error
		if errors.As(err, &apiErr) && apiErr.Code == "chat.not.found" {
			return fmt.Errorf(
				"%w: chat_id=%d",
				errs.ErrChatNotFound,
				chatID,
			)
		}
		return fmt.Errorf("send message through MAX API: %w", err)
	}
	return nil
}
