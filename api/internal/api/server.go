package api

import (
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/middleware"
)

type Handler interface {
	Routes() http.Handler
}

type httpHandler struct {
	userHandler           handler.UserHandler
	eventHandler          handler.EventHandler
	botServiceToken       string
	ingestionServiceToken string
	corsAllowedOrigin     string
}

func NewHandler(
	userHandler handler.UserHandler,
	eventHandler handler.EventHandler,
	botServiceToken string,
	ingestionServiceToken string,
	corsAllowedOrigin string,
) Handler {
	return &httpHandler{
		userHandler:           userHandler,
		eventHandler:          eventHandler,
		botServiceToken:       botServiceToken,
		ingestionServiceToken: ingestionServiceToken,
		corsAllowedOrigin:     corsAllowedOrigin,
	}
}

func (h *httpHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	handler.RegisterRoutes(
		mux,
		h.userHandler,
		h.eventHandler,
		h.botServiceToken,
		h.ingestionServiceToken,
	)

	return middleware.CORS(h.corsAllowedOrigin, mux)
}
