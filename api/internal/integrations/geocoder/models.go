package geocoder

type Address struct {
	BuildingID int     `json:"Building_ID"`
	ID         int     `json:"ID"`
	Name       string  `json:"Name"`
	Longitude  float64 `json:"Longitude"`
	Latitude   float64 `json:"Latitude"`
	DistrictID int     `json:"District_ID"`
	District   string  `json:"District"`
	Flat       string  `json:"Flat"`
}

type Geocode struct {
	ID       int       `json:"id"`
	Address  string    `json:"address"`
	Center   []float64 `json:"center"`
	Distance int       `json:"distance"`
}

type EASAddress struct {
	BuildingID int     `json:"Building_ID"`
	ID         int     `json:"ID"`
	Name       string  `json:"Name"`
	Longitude  float64 `json:"Longitude"`
	Latitude   float64 `json:"Latitude"`
}

type Autocomplete struct {
	AddressID    int     `json:"address_id"`
	Name         string  `json:"Name"`
	BuildingID   int     `json:"building_id"`
	BuildingName string  `json:"building_name"`
	Longitude    float64 `json:"Longitude"`
	Latitude     float64 `json:"Latitude"`
}
