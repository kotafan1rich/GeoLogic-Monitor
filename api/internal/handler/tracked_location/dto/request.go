package dto

import "uuid"

type CreateRequest struct {
	BusinessTypeID *uuid.UUID `json:"business_type_id"`
	Address        *string    `json:"address"`
	Lat            *float64   `json:"lat"`
	Lon            *float64   `json:"lon"`
}
