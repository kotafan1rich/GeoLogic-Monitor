package health

import (
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/health/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
)

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) Health(w http.ResponseWriter, _ *http.Request) {
	response.WriteJSON(w, http.StatusOK, dto.HealthResponse{Status: "ok"})
}
