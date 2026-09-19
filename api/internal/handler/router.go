package handler

import (
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/middleware"
)

type UserHandler interface {
	Upsert(w http.ResponseWriter, r *http.Request)
}

func RegisterRoutes(mux *http.ServeMux, userHandler UserHandler, botServiceToken string) {
	mux.Handle(
		"PUT /api/v1/users/me",
		middleware.BotToken(
			botServiceToken,
			middleware.MaxUserID(http.HandlerFunc(userHandler.Upsert)),
		),
	)
}
