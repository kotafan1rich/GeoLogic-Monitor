package http

import (
	"context"
	"net/http"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type WebhookHandler interface {
	Handle(ctx context.Context, update model.Update)
}

func RegisterRoutes(mux *http.ServeMux, maxClient *maxbot.Api, webhookHandler WebhookHandler, WebhookSecret string) {
	mux.Handle("POST /webhook", maxClient.GetHandler(
		webhookHandler.Handle,
		WebhookSecret,
	))
}
