package geocoding

import (
	"context"
	"errors"
	"math"
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
	Reverse(ctx context.Context, geopoint *domain.GeoPoint) (*domain.Address, error)
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

func (h *handler) Address(w http.ResponseWriter, r *http.Request) {
	lat, err := parseFloat(r.URL.Query().Get("lat"), "lat")
	if err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}
	lon, err := parseFloat(r.URL.Query().Get("lon"), "lon")
	if err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}

	geopoint, err := domain.NewGeoPoint(lat, lon)
	if err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}
	address, err := h.service.Reverse(r.Context(), geopoint)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, dto.ToResponse(address))
}

func parseFloat(value, name string) (float64, error) {
	if value == "" {
		return 0, errors.New(name + " is required")
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(parsed) || math.IsInf(parsed, 0) {
		return 0, errors.New("invalid " + name)
	}
	return parsed, nil
}
