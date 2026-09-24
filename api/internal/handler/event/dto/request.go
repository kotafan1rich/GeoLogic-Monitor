package dto

import "time"

type UpsertRequest struct {
	Provider   string     `json:"provider"`
	ExternalID string     `json:"external_id"`
	Lat        *float64   `json:"lat"`
	Lon        *float64   `json:"lon"`
	Date       *time.Time `json:"date"`
	Info       *string    `json:"info"`
}
