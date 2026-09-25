package service

import (
	"context"
	"time"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

type BusinessType interface {
	Upsert(ctx context.Context, infraTypeID uuid.UUID) (*domain.BusinessType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BusinessType, error)
	GetAll(ctx context.Context) ([]domain.BusinessType, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Event interface {
	Upsert(ctx context.Context, provider, externalID string, lat, lon float64, date time.Time, info *string) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetUnnotifiedByPeriod(ctx context.Context, from, to time.Time) ([]*domain.Event, error)
	GetUnnotifiedNear(ctx context.Context, geoPoint *domain.GeoPoint, radius *uint16, from, to *time.Time) ([]*domain.Event, error)
	MarkNotified(ctx context.Context, id uuid.UUID) error
}

type Geocoding interface {
	Suggest(ctx context.Context, query string) ([]domain.Address, error)
	Reverse(ctx context.Context, geoPoint *domain.GeoPoint) (*domain.Address, error)
}

type InfraObj interface {
	Upsert(ctx context.Context, externalID string, typeID uuid.UUID, lat, lon float64, address string, name *string) (*domain.InfraObject, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraObject, error)
	Near(ctx context.Context, geoPoint *domain.GeoPoint) ([]*domain.InfraObject, error)
}

type InfraType interface {
	Upsert(ctx context.Context, slug, name string, weight float64, maxRadius uint16) (*domain.InfraType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraType, error)
}

type OSRM interface {
	WalkingDistances(ctx context.Context, src *domain.GeoPoint, dst []*domain.GeoPoint) ([]*float64, error)
}

type Rating interface {
	Calculate(ctx context.Context, features domain.LocationFeatures) (*domain.CalculatedRating, error)
	Create(ctx context.Context, trackedLocationID uuid.UUID, rating *domain.CalculatedRating) (*domain.LocationRating, error)
}

type RatingHistory interface {
	GetHistory(ctx context.Context, maxUserID int64, trackedLocationID uuid.UUID, months uint) (*domain.LocationRatingHistory, error)
}

type TrackedLocation interface {
	Create(ctx context.Context, userID, businessTypeID uuid.UUID, name, address string, lat, lon float64) (*domain.TrackedLocationRating, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TrackedLocation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.TrackedLocation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteForUser(ctx context.Context, id, userID uuid.UUID) error
	GetAllForMonitoring(ctx context.Context) ([]domain.MonitoringLocation, error)
	Recalculate(ctx context.Context, id uuid.UUID) error
	RecalculateAll(ctx context.Context) error
}

type User interface {
	Upsert(ctx context.Context, maxUserID, maxChatID int64) (*domain.User, error)
	GetByMaxUserID(ctx context.Context, maxUserID int64) (*domain.User, error)
}
