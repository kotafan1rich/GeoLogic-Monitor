package dto

type Notification struct {
	Type            string          `json:"type"`
	MaxChatID       int64           `json:"max_chat_id"`
	TrackedLocation TrackedLocation `json:"tracked_location"`
	Subject         Subject         `json:"subject"`
	Route           Route           `json:"route"`
	Assessment      Assessment      `json:"assessment"`
}

type TrackedLocation struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	BusinessTypeID string `json:"business_type_id"`
}

type Subject struct {
	ExternalID string `json:"external_id"`
	Title      string `json:"title"`
	Category   string `json:"category,omitempty"`
	Address    string `json:"address,omitempty"`
	OpenedAt   string `json:"opened_at,omitempty"`
	Date       string `json:"date,omitempty"`
}

type Route struct {
	DistanceMeters  float64 `json:"distance_meters"`
	DurationSeconds float64 `json:"duration_seconds"`
}

type Assessment struct {
	Reasons         []string `json:"reasons"`
	Recommendations []string `json:"recommendations"`
}
