package notification

import "github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"

type Notification struct {
	Type            string
	MaxChatID       int64
	TrackedLocation domain.TrackedLocation
	Subject         Subject
	Route           Route
	Assessment      Assessment
}
