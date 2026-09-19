package dto

import "uuid"

type UpsertTypeRequest struct {
	Slug      *string  `json:"slug"`
	Name      *string  `json:"name"`
	Weight    *float64 `json:"weight"`
	MaxRadius *uint16  `json:"max_radius"`
}

type UpsertObjectRequest struct {
	ID      *uuid.UUID `json:"id"`
	TypeID  *uuid.UUID `json:"type_id"`
	Lat     *float64   `json:"lat"`
	Lon     *float64   `json:"lon"`
	Address *string    `json:"address"`
	Name    *string    `json:"name"`
}
