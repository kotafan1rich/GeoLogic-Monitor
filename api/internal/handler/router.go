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

type GeocodingHandler interface {
	Suggest(w http.ResponseWriter, r *http.Request)
}

type TrackedLocationHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetMine(w http.ResponseWriter, r *http.Request)
	DeleteMine(w http.ResponseWriter, r *http.Request)
	GetAllForMonitoring(w http.ResponseWriter, r *http.Request)
}

type BusinessTypeHandler interface {
	Upsert(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type InfraHandler interface {
	UpsertType(w http.ResponseWriter, r *http.Request)
	UpsertObject(w http.ResponseWriter, r *http.Request)
	GetObjectByID(w http.ResponseWriter, r *http.Request)
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
	geocodingHandler GeocodingHandler,
	trackedLocationHandler TrackedLocationHandler,
	businessTypeHandler BusinessTypeHandler,
	infraHandler InfraHandler,
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
	userRoutes := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{"GET /api/v1/geocoding/suggestions", geocodingHandler.Suggest},
		{"POST /api/v1/tracked-locations", trackedLocationHandler.Create},
		{"GET /api/v1/tracked-locations", trackedLocationHandler.GetMine},
		{"DELETE /api/v1/tracked-locations/{id}", trackedLocationHandler.DeleteMine},
	}
	for _, route := range userRoutes {
		mux.Handle(
			route.pattern,
			middleware.BotToken(
				botServiceToken,
				middleware.MaxUserID(route.handler),
			),
		)
	}
	mux.Handle(
		"GET /internal/v1/tracked-locations",
		middleware.IngestionToken(
			ingestionServiceToken,
			http.HandlerFunc(trackedLocationHandler.GetAllForMonitoring),
		),
	)
	mux.Handle(
		"GET /api/v1/business-types",
		middleware.BotToken(
			botServiceToken,
			middleware.MaxUserID(http.HandlerFunc(businessTypeHandler.GetAll)),
		),
	)
	mux.Handle(
		"PUT /internal/v1/business-types",
		middleware.IngestionToken(
			ingestionServiceToken,
			http.HandlerFunc(businessTypeHandler.Upsert),
		),
	)
	infraRoutes := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{"PUT /internal/v1/infra-types", infraHandler.UpsertType},
		{"PUT /internal/v1/infra", infraHandler.UpsertObject},
		{"GET /internal/v1/infra/{id}", infraHandler.GetObjectByID},
	}
	for _, route := range infraRoutes {
		mux.Handle(
			route.pattern,
			middleware.IngestionToken(ingestionServiceToken, route.handler),
		)
	}

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
