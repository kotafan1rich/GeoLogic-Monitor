package competitor

import (
	"time"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/repository/postgresql/geo"
)

type Competitor struct {
	TrackedLocationID uuid.UUID
	ExternalID        string
	Name              string
	TypeID            uuid.UUID
	Address           string
	Location          geo.GeoPoint
	OpenedAt          time.Time
	NotifiedAt        time.Time
	ExpiresAt         time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
