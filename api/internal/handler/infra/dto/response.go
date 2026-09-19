package dto

import "uuid"

type InfraTypeResponse struct {
	ID        uuid.UUID `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Weight    float64   `json:"weight"`
	MaxRadius uint16    `json:"max_radius"`
}

type InfraObjectResponse struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"external_id"`
	TypeID     uuid.UUID `json:"type_id"`
	Lat        float64   `json:"lat"`
	Lon        float64   `json:"lon"`
	Address    string    `json:"address"`
	Name       *string   `json:"name"`
}
