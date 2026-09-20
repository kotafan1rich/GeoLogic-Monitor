package digitalspb

type SpbClassifGateResponse[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type EgsGateResponseV1[T any] struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	List   []T    `json:"list"`
}

type EgsGateResponseV2[T any] struct {
	Count      int `json:"count"`
	TotalPages int `json:"totalPages"`
	Data       []T `json:"data"`
}

type YazzhGateResponse[T any] struct {
	Count int `json:"count"`
	Data  []T `json:"data"`
}

type (
	RailwayStation struct {
		Number               int       `json:"number"`
		Name                 string    `json:"name"`
		AbbreviatedName      string    `json:"abbreviated_name"`
		Address              string    `json:"address"`
		District             string    `json:"district"`
		NearestSubwayStation string    `json:"nearest_subway_station"`
		Chief                string    `json:"chief"`
		WebSite              string    `json:"web_site"`
		Phone                string    `json:"phone"`
		Fax                  string    `json:"fax"`
		INN                  string    `json:"inn"`
		OGRN                 string    `json:"ogrn"`
		Note                 string    `json:"note"`
		Coordinates          []float64 `json:"coordinates"`
	}

	Subway struct {
		Title   string  `json:"title"`
		Address string  `json:"address"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
		Source  string  `json:"source"`
	}
)

type (
	Restaurant struct {
		OID           int       `json:"oid"`
		Name          string    `json:"name"`
		NameEn        string    `json:"name_en"`
		AddressManual string    `json:"address_manual"`
		Phone         string    `json:"phone"`
		WWW           string    `json:"www"`
		Email         string    `json:"email"`
		Kitchen       string    `json:"kitchen"`
		ForDisabled   string    `json:"for_disabled"`
		Coord         []float64 `json:"coord"`
	}

	Hotel struct {
		OID           int       `json:"oid"`
		Name          string    `json:"name"`
		Class         string    `json:"class"`
		HotelRooms    int       `json:"hotel_rooms"`
		HotelsBads    int       `json:"hotels_bads"`
		AddressManual string    `json:"address_manual"`
		District      string    `json:"district"`
		Phone         string    `json:"phone"`
		Fax           string    `json:"fax"`
		WWW           string    `json:"www"`
		Email         string    `json:"email"`
		Description   string    `json:"description"`
		Services      string    `json:"services"`
		Coord         []float64 `json:"coord"`
	}
)

type (
	Cinema struct {
		Number               int    `json:"number"`
		Name                 string `json:"name"`
		AbbreviatedName      string `json:"abbreviated_name"`
		Address              string `json:"address"`
		LegalAddress         string `json:"legal_address"`
		District             string `json:"district"`
		NearestSubwayStation string `json:"nearest_subway_station"`
		Chief                string `json:"chief"`
		Phone                string `json:"phone"`
		Fax                  string `json:"fax"`
		Mode                 string `json:"mode"`
		Email                string `json:"email"`
		MO                   string `json:"mo"`
		INN                  int64  `json:"inn"`
		OGRN                 int64  `json:"ogrn"`
		Note                 string `json:"note"`
	}

	Museum struct {
		OID           int       `json:"oid"`
		Name          string    `json:"name"`
		NameEn        string    `json:"name_en"`
		Type          string    `json:"type"`
		Country       string    `json:"country"`
		AddressManual string    `json:"address_manual"`
		District      string    `json:"district"`
		Phone         string    `json:"phone"`
		WWW           string    `json:"www"`
		Email         string    `json:"email"`
		Description   string    `json:"description"`
		DescriptionEn string    `json:"description_en"`
		WorkTime      string    `json:"work_time"`
		ForDisabled   string    `json:"for_disabled"`
		Coord         []float64 `json:"coord"`
		OGRN          string    `json:"ogrn"`
		INN           string    `json:"inn"`
	}

	ExhibitionHall struct {
		OID           int       `json:"oid"`
		Name          string    `json:"name"`
		NameEn        string    `json:"name_en"`
		Type          string    `json:"type"`
		AddressManual string    `json:"address_manual"`
		Phone         string    `json:"phone"`
		WWW           string    `json:"www"`
		Email         string    `json:"email"`
		Description   string    `json:"description"`
		DescriptionEn string    `json:"description_en"`
		ForDisabled   string    `json:"for_disabled"`
		Coord         []float64 `json:"coord"`
	}

	Theatre struct {
		OID           int       `json:"oid"`
		Name          string    `json:"name"`
		NameEn        string    `json:"name_en"`
		Type          string    `json:"type"`
		AddressManual string    `json:"address_manual"`
		Phone         string    `json:"phone"`
		WWW           string    `json:"www"`
		Email         string    `json:"email"`
		ForDisabled   string    `json:"for_disabled"`
		OGRN          string    `json:"ogrn"`
		INN           string    `json:"inn"`
		Coord         []float64 `json:"coord"`
	}

	Pharmacy struct {
		Number          int       `json:"number"`
		Name            string    `json:"name"`
		NetworkPharmacy string    `json:"network_pharmacy"`
		AbbreviatedName string    `json:"abbreviated_name"`
		Address         string    `json:"address"`
		Coordinates     []float64 `json:"coordinates"`
		District        string    `json:"district"`
		Phone           string    `json:"phone"`
		Mode            string    `json:"mode"`
		OGRN            int64     `json:"ogrn"`
	}

	VetClinic struct {
		Number       int    `json:"number"`
		Name         string `json:"name"`
		DistrictCity string `json:"district_city"`
		Address      string `json:"address"`
		Phone        string `json:"phone"`
		Notes        string `json:"notes"`
	}

	KidsPlace struct {
		ID             int       `json:"id"`
		CategoriesName string    `json:"categories_name"`
		Title          string    `json:"title"`
		Coordinates    []float64 `json:"coordinates"`
	}
)

type (
	StreetPerformance struct {
		ID               int       `json:"id"`
		Title            string    `json:"title"`
		Description      string    `json:"description"`
		DescriptionShort string    `json:"description_short"`
		Note             string    `json:"note"`
		Categories       []string  `json:"categories"`
		Canceled         bool      `json:"canceled"`
		Link             string    `json:"link"`
		StartDate        string    `json:"start_date"`
		EndDate          string    `json:"end_date"`
		Age              string    `json:"age"`
		Source           string    `json:"source"`
		Address          string    `json:"address"`
		Coordinates      []float64 `json:"coordinates"`
		Organizer        string    `json:"organizer"`
		Phone            string    `json:"phone"`
		Format           string    `json:"format:"`
	}

	CultureEvent struct {
		ID                 string                             `json:"ID"`
		Name               string                             `json:"NAME"`
		PreviewText        string                             `json:"PREVIEW_TEXT"`
		DetailText         string                             `json:"DETAIL_TEXT"`
		PreviewPicture     string                             `json:"PREVIEW_PICTURE"`
		PlaceID            string                             `json:"PLACE_ID"`
		PlaceName          string                             `json:"PLACE_NAME"`
		PlaceAddress       string                             `json:"PLACE_ADDRESS"`
		PlaceMap           string                             `json:"PLACE_MAP"`
		Start              string                             `json:"START"`
		Finish             string                             `json:"FINISH"`
		Sessions           []string                           `json:"SESSIONS"`
		Address            string                             `json:"ADDRESS"`
		Map                string                             `json:"MAP"`
		Site               string                             `json:"SITE"`
		Email              string                             `json:"EMAIL"`
		Organizer          string                             `json:"ORGANIZER"`
		Calendar           string                             `json:"CALENDAR"`
		AgeLimitID         string                             `json:"AGE_LIMIT_ID"`
		CategoryID         []string                           `json:"CATEGORY_ID"`
		CalendarOfEventsID []string                           `json:"CALENDAR_OF_EVENTS_ID"`
		Translations       map[string]CultureEventTranslation `json:"TRANSLATIONS"`
	}

	CultureEventTranslation struct {
		PlaceName    string `json:"PLACE_NAME"`
		PlaceAddress string `json:"PLACE_ADDRESS"`
	}
)

type (
	MSPProperty struct {
		ObjID           int64  `json:"obj_id"`
		District        string `json:"district"`
		Address         string `json:"address"`
		Type            string `json:"type"`
		CadastralNumber string `json:"cadastral_number"`
		Area            string `json:"area"`
		Status          string `json:"status"`
		Document        string `json:"document"`
	}
)
