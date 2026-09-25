package trackedlocation

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/request"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/tracked_location/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/middleware"
)

type UserService interface {
	GetByMaxUserID(ctx context.Context, maxUserID int64) (*domain.User, error)
}

type TrackedLocationService interface {
	Create(
		ctx context.Context,
		userID uuid.UUID,
		businessTypeID uuid.UUID,
		name string,
		address string,
		lat float64,
		lon float64,
	) (*domain.TrackedLocationRating, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.TrackedLocation, error)
	DeleteForUser(ctx context.Context, id, userID uuid.UUID) error
	GetAllForMonitoring(ctx context.Context) ([]domain.MonitoringLocation, error)
}

type RatingHistoryService interface {
	GetHistory(
		ctx context.Context,
		maxUserID int64,
		trackedLocationID uuid.UUID,
		months uint,
	) (*domain.LocationRatingHistory, error)
}

type handler struct {
	userService          UserService
	locationService      TrackedLocationService
	ratingHistoryService RatingHistoryService
}

func New(
	userService UserService,
	locationService TrackedLocationService,
	ratingHistoryService RatingHistoryService,
) *handler {
	return &handler{
		userService:          userService,
		locationService:      locationService,
		ratingHistoryService: ratingHistoryService,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateRequest
	if err := request.DecodeJSON(r, &body); err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}
	if body.Name == nil || body.BusinessTypeID == nil || body.Address == nil || body.Lat == nil || body.Lon == nil {
		response.WriteError(w, app.ValidationError(errors.New("name, business_type_id, address, lat and lon are required")))
		return
	}

	user, ok := h.currentUser(w, r)
	if !ok {
		return
	}
	location, err := h.locationService.Create(
		r.Context(),
		user.ID,
		*body.BusinessTypeID,
		*body.Name,
		*body.Address,
		*body.Lat,
		*body.Lon,
	)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, dto.ToCreatedResponse(*location))
}

func (h *handler) GetMine(w http.ResponseWriter, r *http.Request) {
	user, ok := h.currentUser(w, r)
	if !ok {
		return
	}
	locations, err := h.locationService.GetByUserID(r.Context(), user.ID)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponseList(locations))
}

func (h *handler) GetRatingHistory(w http.ResponseWriter, r *http.Request) {
	id, err := request.ParseUUIDPath(r, "id")
	if err != nil {
		response.WriteError(w, app.ValidationError(errors.New("invalid tracked location id")))
		return
	}

	months := uint(3)
	if value := r.URL.Query().Get("months"); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 31)
		if err != nil || parsed < 1 {
			response.WriteError(w, app.ValidationError(errors.New("invalid months")))
			return
		}
		months = uint(parsed)
	}

	maxUserID, ok := middleware.GetMaxUserID(r.Context())
	if !ok {
		response.WriteError(w, app.ErrUnauthorized)
		return
	}

	history, err := h.ratingHistoryService.GetHistory(r.Context(), maxUserID, id, months)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToRatingHistoryResponse(history))
}

func (h *handler) DeleteMine(w http.ResponseWriter, r *http.Request) {
	id, err := request.ParseUUIDPath(r, "id")
	if err != nil {
		response.WriteError(w, app.ValidationError(errors.New("invalid tracked location id")))
		return
	}
	user, ok := h.currentUser(w, r)
	if !ok {
		return
	}
	if err := h.locationService.DeleteForUser(r.Context(), id, user.ID); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) GetAllForMonitoring(w http.ResponseWriter, r *http.Request) {
	locations, err := h.locationService.GetAllForMonitoring(r.Context())
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToMonitoringResponseList(locations))
}

func (h *handler) currentUser(w http.ResponseWriter, r *http.Request) (*domain.User, bool) {
	maxUserID, ok := middleware.GetMaxUserID(r.Context())
	if !ok {
		response.WriteError(w, app.ErrUnauthorized)
		return nil, false
	}
	user, err := h.userService.GetByMaxUserID(r.Context(), maxUserID)
	if err != nil {
		response.WriteServiceError(w, err)
		return nil, false
	}
	return user, true
}
