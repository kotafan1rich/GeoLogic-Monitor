package infra

import (
	"context"
	"errors"
	"net/http"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/infra/dto"
	handlerrequest "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/request"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
)

type InfraTypeService interface {
	Upsert(
		ctx context.Context,
		slug string,
		name string,
		weight float64,
		maxRadius uint16,
	) (*domain.InfraType, error)
}

type InfraService interface {
	Upsert(
		ctx context.Context,
		externalID string,
		typeID uuid.UUID,
		lat float64,
		lng float64,
		address string,
		name *string,
	) (*domain.InfraObject, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraObject, error)
}

type handler struct {
	typeService  InfraTypeService
	infraService InfraService
}

func New(typeService InfraTypeService, infraService InfraService) *handler {
	return &handler{
		typeService:  typeService,
		infraService: infraService,
	}
}

func (h *handler) UpsertType(w http.ResponseWriter, r *http.Request) {
	var request dto.UpsertTypeRequest
	if err := handlerrequest.DecodeJSON(r, &request); err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}
	if request.Slug == nil || request.Name == nil || request.Weight == nil || request.MaxRadius == nil {
		response.WriteError(w, app.ValidationError(errors.New("slug, name, weight and max_radius are required")))
		return
	}

	infraType, err := h.typeService.Upsert(
		r.Context(),
		*request.Slug,
		*request.Name,
		*request.Weight,
		*request.MaxRadius,
	)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.InfraTypeToResponse(*infraType))
}

func (h *handler) UpsertObject(w http.ResponseWriter, r *http.Request) {
	var request dto.UpsertObjectRequest
	if err := handlerrequest.DecodeJSON(r, &request); err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}
	if request.ExternalID == nil || request.TypeID == nil || request.Lat == nil ||
		request.Lon == nil || request.Address == nil {
		response.WriteError(w, app.ValidationError(errors.New("external_id, type_id, lat, lon and address are required")))
		return
	}

	infraObject, err := h.infraService.Upsert(
		r.Context(),
		*request.ExternalID,
		*request.TypeID,
		*request.Lat,
		*request.Lon,
		*request.Address,
		request.Name,
	)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.InfraObjectToResponse(*infraObject))
}

func (h *handler) GetObjectByID(w http.ResponseWriter, r *http.Request) {
	id, err := handlerrequest.ParseUUIDPath(r, "id")
	if err != nil {
		response.WriteError(w, app.ValidationError(errors.New("invalid infra object id")))
		return
	}

	infraObject, err := h.infraService.GetByID(r.Context(), id)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.InfraObjectToResponse(*infraObject))
}
