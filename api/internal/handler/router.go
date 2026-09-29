package handler

import (
	"net/http"
	"time"
)

type HealthHandler interface {
	RegisterRoutes(mux *http.ServeMux)
}

type UserHandler interface {
	RegisterRoutes(mux *http.ServeMux, botServiceToken string)
}

type GeocodingHandler interface {
	RegisterRoutes(mux *http.ServeMux, ingestionServiceToken, maxBotToken string, miniAppInitDataMaxAge time.Duration)
}

type TrackedLocationHandler interface {
	RegisterRoutes(mux *http.ServeMux, ingestionServiceToken, maxBotToken string, miniAppInitDataMaxAge time.Duration)
}

type BusinessTypeHandler interface {
	RegisterRoutes(mux *http.ServeMux, ingestionServiceToken, maxBotToken string, miniAppInitDataMaxAge time.Duration)
}

type InfraHandler interface {
	RegisterRoutes(mux *http.ServeMux, ingestionServiceToken string)
}

type EventHandler interface {
	RegisterRoutes(mux *http.ServeMux, ingestionServiceToken string)
}

type RoutesHandler interface {
	RegisterRoutes(mux *http.ServeMux, ingestionServiceToken string)
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
	routesHandler RoutesHandler,
	maxBotToken string,
	miniAppInitDataMaxAge time.Duration,
	botServiceToken string,
	ingestionServiceToken string,
) {
	healthHandler.RegisterRoutes(mux)
	userHandler.RegisterRoutes(mux, botServiceToken)
	geocodingHandler.RegisterRoutes(mux, ingestionServiceToken, maxBotToken, miniAppInitDataMaxAge)
	trackedLocationHandler.RegisterRoutes(mux, ingestionServiceToken, maxBotToken, miniAppInitDataMaxAge)
	businessTypeHandler.RegisterRoutes(mux, ingestionServiceToken, maxBotToken, miniAppInitDataMaxAge)
	infraHandler.RegisterRoutes(mux, ingestionServiceToken)
	eventHandler.RegisterRoutes(mux, ingestionServiceToken)
	routesHandler.RegisterRoutes(mux, ingestionServiceToken)
}
