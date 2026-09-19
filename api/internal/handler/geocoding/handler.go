package geocoding

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/geocoding/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
)

const defaultLimit = 5

type GeocodingService interface {
	Suggest(ctx context.Context, query string) ([]domain.Address, error)
}

type handler struct {
	service GeocodingService
}

func New(service GeocodingService) *handler {
	return &handler{service: service}
}

func (h *handler) Suggest(w http.ResponseWriter, r *http.Request) {
	limit := defaultLimit
	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > defaultLimit {
			response.WriteError(w, app.ValidationError(errors.New("invalid limit")))
			return
		}
		limit = parsed
	}

	addresses, err := h.service.Suggest(r.Context(), r.URL.Query().Get("query"))
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}
	if len(addresses) > limit {
		addresses = addresses[:limit]
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponseList(addresses))
}
