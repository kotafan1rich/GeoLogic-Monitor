package server

import (
	"net/http"

	handler "github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/http"
	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
)

type Handler interface {
	Routes() http.Handler
}

type httpHandler struct {
	maxClient      *maxbot.Api
	webhookHandler handler.WebhookHandler
	webhookSecret  string
}

func NewHandler(
	maxClient *maxbot.Api,
	webhookHandler handler.WebhookHandler,
	webhookSecret string,
) *httpHandler {
	return &httpHandler{
		maxClient:      maxClient,
		webhookHandler: webhookHandler,
		webhookSecret:  webhookSecret,
	}
}

func (h *httpHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, h.maxClient, h.webhookHandler, h.webhookSecret)

	return mux
}
