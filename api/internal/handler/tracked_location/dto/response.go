package dto

import (
	"time"
	"uuid"
)

type TrackedLocationResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	BusinessTypeID uuid.UUID `json:"business_type_id"`
	Address        string    `json:"address"`
	Lat            float64   `json:"lat"`
	Lon            float64   `json:"lon"`
}

type CreatedTrackedLocationResponse struct {
	TrackedLocationResponse
	Rating             float64   `json:"rating"`
	RatingCalculatedAt time.Time `json:"rating_calculated_at"`
}

type UserTrackedLocationResponse struct {
	TrackedLocationResponse
	Rating             *float64   `json:"rating"`
	RatingCalculatedAt *time.Time `json:"rating_calculated_at"`
}

type RatingHistoryEntryResponse struct {
	Value        float64   `json:"value"`
	CalculatedAt time.Time `json:"calculated_at"`
}

type MonitoringLocationResponse struct {
	TrackedLocationResponse
	MaxChatID int64 `json:"max_chat_id"`
}
