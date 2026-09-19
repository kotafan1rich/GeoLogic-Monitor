package event

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"
	_ "time/tzdata"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/event/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
)

const dateLayout = "2006-01-02"

type EventService interface {
	Upsert(ctx context.Context, provider string, externalID string, lat float64, lng float64, date time.Time, info *string) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetUnnotifiedByPeriod(ctx context.Context, from, to time.Time) ([]domain.Event, error)
	GetUnnotifiedNear(
		ctx context.Context,
		geoPoint *domain.GeoPoint,
		radius *uint16,
		from *time.Time,
		to *time.Time,
	) ([]domain.Event, error)
	MarkNotified(ctx context.Context, id uuid.UUID) error
}

type handler struct {
	service        EventService
	moscowLocation *time.Location
}

func New(service EventService) *handler {
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("event handler: Europe/Moscow timezone is unavailable")
	}

	return &handler{
		service:        service,
		moscowLocation: moscowLocation,
	}
}

func (h *handler) Upsert(w http.ResponseWriter, r *http.Request) {
	var request dto.UpsertRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		response.WriteError(w, app.ValidationError(errors.New("invalid request body")))
		return
	}

	if request.Lat == nil || request.Lon == nil || request.Date == nil {
		response.WriteError(w, app.ValidationError(errors.New("lat, lon and date are required")))
		return
	}

	event, err := h.service.Upsert(
		r.Context(),
		request.Provider,
		request.ExternalID,
		*request.Lat,
		*request.Lon,
		*request.Date,
		request.Info,
	)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponse(*event))
}

func (h *handler) ListUnnotified(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.parseDay(r.URL.Query().Get("day"))
	if err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}

	events, err := h.service.GetUnnotifiedByPeriod(r.Context(), from, to)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponseList(events))
}

func (h *handler) ListUnnotifiedNear(w http.ResponseWriter, r *http.Request) {
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

	geoPoint, err := domain.NewGeoPoint(lat, lon)
	if err != nil {
		response.WriteError(w, app.ValidationError(err))
		return
	}

	var radius *uint16
	if value := r.URL.Query().Get("radius"); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			response.WriteError(w, app.ValidationError(errors.New("invalid radius")))
			return
		}
		value := uint16(parsed)
		radius = &value
	}

	var from, to *time.Time
	if day := r.URL.Query().Get("day"); day != "" {
		parsedFrom, parsedTo, err := h.parseDay(day)
		if err != nil {
			response.WriteError(w, app.ValidationError(err))
			return
		}
		from, to = &parsedFrom, &parsedTo
	}

	events, err := h.service.GetUnnotifiedNear(r.Context(), geoPoint, radius, from, to)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponseList(events))
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, app.ValidationError(errors.New("invalid event id")))
		return
	}

	event, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, dto.ToResponse(*event))
}

func (h *handler) MarkNotified(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, app.ValidationError(errors.New("invalid event id")))
		return
	}

	if err := h.service.MarkNotified(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) parseDay(value string) (time.Time, time.Time, error) {
	if value == "" {
		return time.Time{}, time.Time{}, errors.New("day is required")
	}

	from, err := time.ParseInLocation(dateLayout, value, h.moscowLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid day")
	}

	return from, from.AddDate(0, 0, 1), nil
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

func writeServiceError(w http.ResponseWriter, err error) {
	if serviceError, ok := errors.AsType[*app.Error](err); ok {
		response.WriteError(w, serviceError)
	} else {
		response.WriteError(w, app.ErrInternal)
	}
}
