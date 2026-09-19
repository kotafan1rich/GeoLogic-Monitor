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
	healthHandler          handler.HealthHandler
	userHandler            handler.UserHandler
	geocodingHandler       handler.GeocodingHandler
	trackedLocationHandler handler.TrackedLocationHandler
	businessTypeHandler    handler.BusinessTypeHandler
	infraHandler           handler.InfraHandler
	eventHandler           handler.EventHandler
	botServiceToken        string
	ingestionServiceToken  string
	corsAllowedOrigin      string
}

func NewHandler(
	healthHandler handler.HealthHandler,
	userHandler handler.UserHandler,
	geocodingHandler handler.GeocodingHandler,
	trackedLocationHandler handler.TrackedLocationHandler,
	businessTypeHandler handler.BusinessTypeHandler,
	infraHandler handler.InfraHandler,
	eventHandler handler.EventHandler,
	botServiceToken string,
	ingestionServiceToken string,
	corsAllowedOrigin string,
) Handler {
	return &httpHandler{
		healthHandler:          healthHandler,
		userHandler:            userHandler,
		geocodingHandler:       geocodingHandler,
		trackedLocationHandler: trackedLocationHandler,
		businessTypeHandler:    businessTypeHandler,
		infraHandler:           infraHandler,
		eventHandler:           eventHandler,
		botServiceToken:        botServiceToken,
		ingestionServiceToken:  ingestionServiceToken,
		corsAllowedOrigin:      corsAllowedOrigin,
	}
}

func (h *httpHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	handler.RegisterRoutes(
		mux,
		h.healthHandler,
		h.userHandler,
		h.geocodingHandler,
		h.trackedLocationHandler,
		h.businessTypeHandler,
		h.infraHandler,
		h.eventHandler,
		h.botServiceToken,
		h.ingestionServiceToken,
	)

	return middleware.CORS(h.corsAllowedOrigin, mux)
}
