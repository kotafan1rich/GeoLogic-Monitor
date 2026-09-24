package model

import (
	"time"
	"uuid"
)

type Base struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
