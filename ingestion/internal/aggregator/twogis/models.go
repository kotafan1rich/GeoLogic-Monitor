package twogis

type (
	Response struct {
		Meta   Meta   `json:"meta"`
		Result Result `json:"result"`
	}

	Meta struct {
		APIVersion string     `json:"api_version"`
		Code       int        `json:"code"`
		Error      *MetaError `json:"error,omitempty"`
	}

	MetaError struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}

	Result struct {
		Items []Item `json:"items"`
		Total int    `json:"total"`
	}

	Item struct {
		ID          string   `json:"id"`
		Type        string   `json:"type"`
		Name        string   `json:"name"`
		AddressName string   `json:"address_name"`
		Point       *Point   `json:"point"`
		Address     *Address `json:"address"`
		Rubrics     []Rubric `json:"rubrics"`
		Org         *Org     `json:"org"`
	}

	Point struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	Address struct {
		BuildingID   string             `json:"building_id"`
		BuildingName string             `json:"building_name"`
		Postcode     string             `json:"postcode"`
		Components   []AddressComponent `json:"components"`
	}

	AddressComponent struct {
		Type   string `json:"type"`
		Street string `json:"street"`
		Number string `json:"number"`
	}

	Rubric struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Alias    string `json:"alias"`
		ParentID string `json:"parent_id"`
	}

	Org struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		BranchCount int    `json:"branch_count"`
	}
)
