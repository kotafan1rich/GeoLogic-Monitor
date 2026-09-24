package geocoder

type suggestRequest struct {
	Query         string            `json:"query"`
	Count         int               `json:"count"`
	Locations     []suggestLocation `json:"locations"`
	RestrictValue bool              `json:"restrict_value"`
}

type suggestLocation struct {
	City string `json:"city"`
}

type geolocateRequest struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Count int     `json:"count"`
}

type suggestionsResponse struct {
	Suggestions []suggestion `json:"suggestions"`
}

type suggestion struct {
	Value string         `json:"value"`
	Data  suggestionData `json:"data"`
}

type suggestionData struct {
	GeoLat float64 `json:"geo_lat,string"`
	GeoLon float64 `json:"geo_lon,string"`
}

type Autocomplete struct {
	Name         string
	BuildingName string
	Longitude    float64
	Latitude     float64
}

type Geocode struct {
	Address string
	Center  GeocodeCenter
}

type GeocodeCenter struct {
	X float64
	Y float64
}
