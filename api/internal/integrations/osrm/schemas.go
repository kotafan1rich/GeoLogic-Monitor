package osrm

type tableResponse struct {
	Code      string       `json:"code"`
	Distances [][]*float64 `json:"distances"`
}
