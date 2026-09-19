package api

import (
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler"
)

type Handler interface {
	Routes() http.Handler
}

type httpHandler struct {
	userHandler     handler.UserHandler
	botServiceToken string
}

func NewHandler(userHandler handler.UserHandler, botServiceToken string) Handler {
	return &httpHandler{
		userHandler:     userHandler,
		botServiceToken: botServiceToken,
	}
}

func (h *httpHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, h.userHandler, h.botServiceToken)

	return mux
}
