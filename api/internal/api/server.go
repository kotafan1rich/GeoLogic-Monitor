package api

import (
	"net/http"
)

type Handler interface {
	Routes() http.Handler
}

type httpHandler struct {
}

func NewHandler() Handler {
	return &httpHandler{}
}

func (h *httpHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	return mux
}
