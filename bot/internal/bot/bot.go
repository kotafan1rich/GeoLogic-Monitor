package bot

import (
	"net/http"

	handler "github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/http"
)

type Handler interface {
	Routes() http.Handler
}

type httpHandler struct{}

func NewHandler() *httpHandler {
	return &httpHandler{}
}

func (h *httpHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return mux
}
