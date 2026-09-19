package dto

import "uuid"

type InfraTypeResponse struct {
	ID        uuid.UUID `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Weight    float64   `json:"weight"`
	MaxRadius uint16    `json:"max_radius"`
}

type BusinessTypeResponse struct {
	ID          uuid.UUID         `json:"id"`
	InfraTypeID uuid.UUID         `json:"infra_type_id"`
	InfraType   InfraTypeResponse `json:"infra_type"`
}
