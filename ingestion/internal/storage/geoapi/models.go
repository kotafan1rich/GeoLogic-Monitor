package geoapi

import (
	"time"

	"github.com/google/uuid"
)

type (
	AddressComponent struct {
		Address string  `json:"address"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}

	GeoPoint struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	InfraTypeInput struct {
		Slug      string   `json:"slug"`
		Name      string   `json:"name"`
		Weight    int      `json:"weight"`
		MaxRadius int      `json:"max_radius"`
		Query     string   `json:"-"`
		Rubrics   []string `json:"-"`
	}

	InfraType struct {
		ID        string `json:"id"`
		Slug      string `json:"slug"`
		Name      string `json:"name"`
		Weight    int    `json:"weight"`
		MaxRadius int    `json:"max_radius"`
	}

	InfraObjectInput struct {
		ExternalID string  `json:"external_id"`
		TypeID     string  `json:"type_id"`
		Name       *string `json:"name,omitempty"`
		Address    string  `json:"address"`
		Lat        float64 `json:"lat"`
		Lon        float64 `json:"lon"`
	}

	InfraObject struct {
		ID         string  `json:"id"`
		ExternalID string  `json:"external_id"`
		TypeID     string  `json:"type_id"`
		Name       *string `json:"name,omitempty"`
		Address    string  `json:"address"`
		Lat        float64 `json:"lat"`
		Lon        float64 `json:"lon"`
	}

	EventInput struct {
		Provider   string    `json:"provider"`
		ExternalID string    `json:"external_id"`
		Lat        float64   `json:"lat"`
		Lon        float64   `json:"lon"`
		Date       time.Time `json:"date"`
		Info       *string   `json:"info,omitempty"`
	}

	Event struct {
		ID         string    `json:"id"`
		Provider   string    `json:"provider"`
		ExternalID string    `json:"external_id"`
		Lat        float64   `json:"lat"`
		Lon        float64   `json:"lon"`
		Date       time.Time `json:"date"`
		Info       *string   `json:"info,omitempty"`
		NotifiedAt time.Time `json:"notified_at"`
	}

	InfraTypeID struct {
		InfraTypeID string `json:"infra_type_id"`
	}

	BusinessType struct {
		ID          string    `json:"id"`
		InfraTypeID string    `json:"infra_type_id"`
		InfraType   InfraType `json:"infra_type"`
	}

	TrackedLocation struct {
		ID             uuid.UUID `json:"id"`
		Name           string    `json:"name"`
		BusinessTypeID uuid.UUID `json:"business_type_id"`
		Address        string    `json:"address"`
		Lat            float64   `json:"lat"`
		Lon            float64   `json:"lon"`
		MaxChatID      int64     `json:"max_chat_id"`
	}

	WalkingDistancesRequest struct {
		Source       GeoPoint   `json:"source"`
		Destinations []GeoPoint `json:"destinations"`
	}

	WalkingDistancesResponse struct {
		Distances []*float64 `json:"distances"`
	}
)
