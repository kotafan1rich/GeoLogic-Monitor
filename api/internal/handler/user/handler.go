package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	handlerrequest "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/request"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/user/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/middleware"
)

type UserService interface {
	Upsert(ctx context.Context, maxUserID int64, maxChatID int64) (*domain.User, error)
}

type handler struct {
	service UserService
}

func New(service UserService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Upsert(w http.ResponseWriter, r *http.Request) {
	var request dto.UpsertRequest
	if err := handlerrequest.DecodeJSON(r, &request); err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}

	if request.MaxChatID == nil {
		response.WriteError(w, app.ValidationError(errors.New("max_chat_id is required")))
		return
	}

	ctx := r.Context()

	maxUserID, ok := middleware.GetMaxUserID(ctx)
	if !ok {
		response.WriteError(w, app.ErrUnauthorized)
		return
	}

	user, err := h.service.Upsert(ctx, maxUserID, *request.MaxChatID)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponse(*user))
}
