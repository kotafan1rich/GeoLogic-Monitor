package max

import (
	"context"
	"fmt"

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
		return fmt.Errorf("send message through MAX API: %w", err)
	}
	return nil
}
