package handler

import (
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/middleware"
)

type HealthHandler interface {
	Health(w http.ResponseWriter, r *http.Request)
}

type UserHandler interface {
	Upsert(w http.ResponseWriter, r *http.Request)
}

type EventHandler interface {
	Upsert(w http.ResponseWriter, r *http.Request)
	ListUnnotified(w http.ResponseWriter, r *http.Request)
	ListUnnotifiedNear(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	MarkNotified(w http.ResponseWriter, r *http.Request)
}

func RegisterRoutes(
	mux *http.ServeMux,
	healthHandler HealthHandler,
	userHandler UserHandler,
	eventHandler EventHandler,
	botServiceToken string,
	ingestionServiceToken string,
) {
	mux.HandleFunc("GET /health", healthHandler.Health)

	mux.Handle(
		"PUT /api/v1/users/me",
		middleware.BotToken(
			botServiceToken,
			middleware.MaxUserID(http.HandlerFunc(userHandler.Upsert)),
		),
	)

	eventRoutes := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{"PUT /internal/v1/events", eventHandler.Upsert},
		{"GET /internal/v1/events", eventHandler.ListUnnotified},
		{"GET /internal/v1/events/near", eventHandler.ListUnnotifiedNear},
		{"GET /internal/v1/events/{id}", eventHandler.GetByID},
		{"PUT /internal/v1/events/{id}/notified", eventHandler.MarkNotified},
	}
	for _, route := range eventRoutes {
		mux.Handle(
			route.pattern,
			middleware.IngestionToken(ingestionServiceToken, route.handler),
		)
	}
}
