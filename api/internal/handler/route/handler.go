package route

import (
	"context"
	"errors"
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/request"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/route/dto"
)

const maxDestinations = 99

type OSRMService interface {
	WalkingDistances(
		ctx context.Context,
		src *domain.GeoPoint,
		dst []*domain.GeoPoint,
	) ([]*float64, error)
}

type handler struct {
	service OSRMService
}

func New(service OSRMService) *handler {
	return &handler{service: service}
}

func (h *handler) WalkingDistances(w http.ResponseWriter, r *http.Request) {
	var body dto.WalkingDistancesRequest
	if err := request.DecodeJSON(r, &body); err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}

	if body.Source == nil {
		response.WriteError(w, app.ValidationError(errors.New("source is required")))
		return
	}
	if len(body.Destinations) == 0 || len(body.Destinations) > maxDestinations {
		response.WriteError(w, app.ValidationError(errors.New("destinations must contain from 1 to 99 points")))
		return
	}

	source, err := dto.ToGeoPoint(*body.Source)
	if err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}

	destinations := make([]*domain.GeoPoint, 0, len(body.Destinations))
	for _, destination := range body.Destinations {
		point, err := dto.ToGeoPoint(destination)
		if err != nil {
			response.WriteError(w, app.ValidationError(err))
			return
		}
		destinations = append(destinations, point)
	}

	distances, err := h.service.WalkingDistances(r.Context(), source, destinations)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponse(distances))
}
