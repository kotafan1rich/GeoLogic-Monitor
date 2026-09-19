package dto

import (
	"time"
	"uuid"
)

type EventResponse struct {
	ID         uuid.UUID  `json:"id"`
	Provider   string     `json:"provider"`
	ExternalID string     `json:"external_id"`
	Lat        float64    `json:"lat"`
	Lon        float64    `json:"lon"`
	Date       time.Time  `json:"date"`
	Info       *string    `json:"info"`
	NotifiedAt *time.Time `json:"notified_at"`
}
