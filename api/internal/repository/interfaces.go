package repository

import (
	"context"
	"time"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
)

type BusinessType interface {
	Upsert(ctx context.Context, businessType *domain.BusinessType) (*domain.BusinessType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BusinessType, error)
	GetAll(ctx context.Context) ([]domain.BusinessType, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Event interface {
	Upsert(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetUnnotifiedByPeriod(ctx context.Context, from time.Time, to time.Time) ([]*domain.Event, error)
	GetUnnotifiedNear(ctx context.Context, geoPoint *domain.GeoPoint, radius uint16, from *time.Time, to *time.Time) ([]*domain.Event, error)
	MarkNotified(ctx context.Context, id uuid.UUID) error
}

type InfraObject interface {
	Upsert(ctx context.Context, infraObject *domain.InfraObject) (*domain.InfraObject, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraObject, error)
	Near(ctx context.Context, geoPoint *domain.GeoPoint) ([]*domain.InfraObject, error)
}

type InfraType interface {
	Upsert(ctx context.Context, infraType *domain.InfraType) (*domain.InfraType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.InfraType, error)
}

type Rating interface {
	Create(ctx context.Context, locationRating *domain.LocationRating) (*domain.LocationRating, error)
	CreateMany(ctx context.Context, locationRatings []domain.LocationRating) ([]domain.LocationRating, error)
	GetHistory(ctx context.Context, trackedLocationID uuid.UUID, months uint) (*domain.LocationRatingHistory, error)
}

type TrackedLocation interface {
	Create(ctx context.Context, location *domain.TrackedLocation) (*domain.TrackedLocation, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TrackedLocation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.TrackedLocation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllForMonitoring(ctx context.Context) ([]domain.MonitoringLocation, error)
}

type User interface {
	Upsert(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByMaxUserID(ctx context.Context, maxUserID int64) (*domain.User, error)
}

type Geocoder interface {
	Suggestions(ctx context.Context, query string) ([]domain.Address, error)
	Reverse(ctx context.Context, lat, lon float64) (*domain.Address, error)
}

type OSRM interface {
	WalkingDistances(ctx context.Context, src *domain.GeoPoint, dst []*domain.GeoPoint) ([]*float64, error)
}
