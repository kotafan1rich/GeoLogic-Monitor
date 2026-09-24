package api

import (
	"net/http"
	"time"

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
	routesHandler          handler.RoutesHandler
	maxBotToken            string
	miniAppInitDataMaxAge  time.Duration
	botServiceToken        string
	ingestionServiceToken  string
	corsAllowedOrigin      string
	docsDir                string
}

func NewHandler(
	healthHandler handler.HealthHandler,
	userHandler handler.UserHandler,
	geocodingHandler handler.GeocodingHandler,
	trackedLocationHandler handler.TrackedLocationHandler,
	businessTypeHandler handler.BusinessTypeHandler,
	infraHandler handler.InfraHandler,
	eventHandler handler.EventHandler,
	routesHandler handler.RoutesHandler,
	maxBotToken string,
	miniAppInitDataMaxAge time.Duration,
	botServiceToken string,
	ingestionServiceToken string,
	corsAllowedOrigin string,
	docsDir string,
) Handler {
	return &httpHandler{
		healthHandler:          healthHandler,
		userHandler:            userHandler,
		geocodingHandler:       geocodingHandler,
		trackedLocationHandler: trackedLocationHandler,
		businessTypeHandler:    businessTypeHandler,
		infraHandler:           infraHandler,
		eventHandler:           eventHandler,
		routesHandler:          routesHandler,
		maxBotToken:            maxBotToken,
		miniAppInitDataMaxAge:  miniAppInitDataMaxAge,
		botServiceToken:        botServiceToken,
		ingestionServiceToken:  ingestionServiceToken,
		corsAllowedOrigin:      corsAllowedOrigin,
		docsDir:                docsDir,
	}
}

func (h *httpHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"GET /docs/",
		http.StripPrefix(
			"/docs/",
			http.FileServer(http.Dir(h.docsDir))),
	)

	handler.RegisterRoutes(
		mux,
		h.healthHandler,
		h.userHandler,
		h.geocodingHandler,
		h.trackedLocationHandler,
		h.businessTypeHandler,
		h.infraHandler,
		h.eventHandler,
		h.routesHandler,
		h.maxBotToken,
		h.miniAppInitDataMaxAge,
		h.botServiceToken,
		h.ingestionServiceToken,
	)

	return middleware.CORS(h.corsAllowedOrigin, mux)
}
