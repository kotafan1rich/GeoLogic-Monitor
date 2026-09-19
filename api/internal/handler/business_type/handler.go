package businesstype

import (
	"context"
	"errors"
	"net/http"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/business_type/dto"
	handlerrequest "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/request"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
)

type BusinessTypeService interface {
	Upsert(ctx context.Context, infraTypeID uuid.UUID) (*domain.BusinessType, error)
	GetAll(ctx context.Context) ([]domain.BusinessType, error)
}

type handler struct {
	service BusinessTypeService
}

func New(service BusinessTypeService) *handler {
	return &handler{service: service}
}

func (h *handler) Upsert(w http.ResponseWriter, r *http.Request) {
	var request dto.UpsertRequest
	if err := handlerrequest.DecodeJSON(r, &request); err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}
	if request.InfraTypeID == nil {
		response.WriteError(w, app.ValidationError(errors.New("infra_type_id is required")))
		return
	}

	businessType, err := h.service.Upsert(r.Context(), *request.InfraTypeID)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponse(*businessType))
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	businessTypes, err := h.service.GetAll(r.Context())
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponseList(businessTypes))
}
