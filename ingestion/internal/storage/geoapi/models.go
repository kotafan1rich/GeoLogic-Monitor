package geoapi

import "time"

type (
	AddressComponent struct {
		Address string  `json:"address"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}

	InfraTypeInput struct {
		Slug      string `json:"slug"`
		Name      string `json:"name"`
		Weight    int    `json:"weight"`
		MaxRadius int    `json:"max_radius"`
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

	Event struct {
		Provider   string    `json:"provider"`
		ExternalID string    `json:"external_id"`
		Lat        float64   `json:"lat"`
		Lon        float64   `json:"lon"`
		Date       time.Time `json:"date"`
		Info       *string   `json:"info,omitempty"`
	}
)
