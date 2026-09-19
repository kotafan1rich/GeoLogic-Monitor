package dto

import "uuid"

type TrackedLocationResponse struct {
	ID             uuid.UUID `json:"id"`
	BusinessTypeID uuid.UUID `json:"business_type_id"`
	Address        string    `json:"address"`
	Lat            float64   `json:"lat"`
	Lon            float64   `json:"lon"`
}

type CreatedTrackedLocationResponse struct {
	TrackedLocationResponse
	Rating float64 `json:"rating"`
}

type MonitoringLocationResponse struct {
	TrackedLocationResponse
	MaxChatID int64 `json:"max_chat_id"`
}
