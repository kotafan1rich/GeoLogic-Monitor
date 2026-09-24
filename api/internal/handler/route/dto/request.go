package dto

type GeoPointRequest struct {
	Lat *float64 `json:"lat"`
	Lon *float64 `json:"lon"`
}

type WalkingDistancesRequest struct {
	Source       *GeoPointRequest  `json:"source"`
	Destinations []GeoPointRequest `json:"destinations"`
}
