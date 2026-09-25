package maps

import (
	"fmt"
	"strings"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const providerPrefix = "maps"

const (
	datasetSubway            = "subway"
	datasetRailwayStation    = "railway_station"
	datasetHotel             = "hotel"
	datasetMuseum            = "museum"
	datasetTheatre           = "theatre"
	datasetCinema            = "cinema"
	datasetExhibitionHall    = "exhibition_hall"
	datasetZoo               = "zoo"
	datasetPark              = "park"
	datasetPlayground        = "playground"
	datasetLibrary           = "library"
	datasetCamp              = "camp"
	datasetEducationCenter   = "education_center"
	datasetGameCenter        = "game_center"
	datasetSportsCenter      = "sports_center"
	datasetSupermarket       = "supermarket"
	datasetMall              = "mall"
	datasetPVZ               = "pvz"
	datasetGrocery           = "grocery"
	datasetCafe              = "cafe"
	datasetCoffee            = "coffee"
	datasetRestaurant        = "restaurant"
	datasetFastfood          = "fastfood"
	datasetBakery            = "bakery"
	datasetAlcohol           = "alcohol"
	datasetBeauty            = "beauty"
	datasetPharmacy          = "pharmacy"
	datasetMedicalCenter     = "medical_center"
	datasetDental            = "dental"
	datasetVetClinic         = "vet_clinic"
	datasetPet               = "pet"
	datasetOptician          = "optician"
	datasetFlowers           = "flowers"
	datasetLaundry           = "laundry"
	datasetCopyshop          = "copyshop"
	datasetWorkshop          = "workshop"
	datasetElectronicsRepair = "electronics_repair"
	datasetCarService        = "car_service"
	datasetCarWash           = "car_wash"
)

func MapToInfraObjects(objs []Element) []geoapi.InfraObjectInput {
	mapped := make([]geoapi.InfraObjectInput, 0, len(objs))

	for _, obj := range objs {
		mappedObj := mapToInfraObject(obj)
		if mappedObj == (geoapi.InfraObjectInput{}) {
			continue
		}
		mapped = append(mapped, mappedObj)
	}

	return mapped
}

func mapToInfraObject(obj Element) geoapi.InfraObjectInput {
	if len(obj.Tags) == 0 {
		return geoapi.InfraObjectInput{}
	}

	ds := dataset(obj.Tags)
	if ds == "" {
		return geoapi.InfraObjectInput{}
	}

	lat, lon := coords(obj)
	if lat == 0 || lon == 0 {
		return geoapi.InfraObjectInput{}
	}

	return geoapi.InfraObjectInput{
		ExternalID: externalID(ds, obj.ID),
		TypeID:     "",
		Address:    address(obj.Tags),
		Lat:        lat,
		Lon:        lon,
		Name:       name(obj.Tags),
	}
}

var osmKeys = []string{"shop", "amenity", "leisure", "tourism", "railway", "craft"}

func dataset(t map[string]string) string {
	for _, key := range osmKeys {
		raw, ok := t[key]
		if !ok {
			continue
		}
		for _, part := range strings.Split(raw, ";") {
			if slug := matchSlug(strings.TrimSpace(part), t); slug != "" {
				return slug
			}
		}
	}
	return ""
}

func matchSlug(v string, t map[string]string) string {
	switch v {
	case "subway_entrance":
		return datasetSubway
	case "station", "halt":
		if t["station"] == "subway" {
			return ""
		}
		return datasetRailwayStation
	case "hotel", "hostel", "guest_house":
		return datasetHotel
	case "museum":
		return datasetMuseum
	case "theatre":
		return datasetTheatre
	case "cinema":
		return datasetCinema
	case "gallery", "arts_centre":
		return datasetExhibitionHall
	case "zoo", "aquarium":
		return datasetZoo
	case "park", "garden":
		return datasetPark
	case "playground":
		return datasetPlayground
	case "library":
		return datasetLibrary
	case "camp_site":
		return datasetCamp
	case "school", "university", "college", "kindergarten":
		return datasetEducationCenter
	case "adult_gaming_centre", "escape_game", "amusement_arcade":
		return datasetGameCenter
	case "fitness_centre", "sports_centre", "sports_hall":
		return datasetSportsCenter
	case "supermarket":
		return datasetSupermarket
	case "mall", "department_store":
		return datasetMall
	case "outpost", "parcel_locker":
		return datasetPVZ
	case "convenience", "greengrocer", "butcher", "seafood", "deli", "farm",
		"dairy", "cheese", "nuts", "health_food", "frozen_food", "spices",
		"honey", "chocolate":
		return datasetGrocery
	case "cafe":
		if strings.Contains(t["cuisine"], "coffee_shop") {
			return datasetCoffee
		}
		return datasetCafe
	case "coffee":
		return datasetCoffee
	case "restaurant", "bar", "pub":
		return datasetRestaurant
	case "fast_food", "food_court":
		return datasetFastfood
	case "bakery", "pastry", "confectionery":
		return datasetBakery
	case "alcohol", "wine", "beverages", "brewing_supplies":
		return datasetAlcohol
	case "beauty", "hairdresser", "massage", "tattoo", "piercing", "spa", "nails":
		return datasetBeauty
	case "pharmacy", "chemist", "medical_supply", "nutrition_supplements",
		"orthopedics", "herbalist":
		return datasetPharmacy
	case "clinic", "doctors":
		return datasetMedicalCenter
	case "dentist":
		return datasetDental
	case "veterinary":
		return datasetVetClinic
	case "pet", "pet_grooming":
		return datasetPet
	case "optician", "hearing_aids":
		return datasetOptician
	case "florist":
		return datasetFlowers
	case "laundry", "dry_cleaning":
		return datasetLaundry
	case "copyshop", "photo", "photo_studio", "printing", "photographer":
		return datasetCopyshop
	case "shoe_repair", "shoemaker", "tailor", "sewing", "watches", "watchmaker",
		"key_cutter", "locksmith", "repair":
		return datasetWorkshop
	case "mobile_phone", "electronics", "computer", "radiotechnics",
		"hifi", "electronics_repair":
		return datasetElectronicsRepair
	case "car_repair", "car_parts", "tyres", "tyres_repair", "motorcycle_repair":
		return datasetCarService
	case "car_wash":
		return datasetCarWash
	}

	return ""
}

func externalID(dataset string, id int64) string {
	return fmt.Sprintf("%s:%s:%d", providerPrefix, dataset, id)
}

func address(t map[string]string) string {
	street, ok := t["addr:street"]
	if !ok {
		return "mock-address"
	}

	houseNumber, ok := t["addr:housenumber"]
	if !ok {
		return street
	}

	return fmt.Sprintf("%s, %s", street, houseNumber)
}

func coords(obj Element) (lat, lon float64) {
	if obj.Lat != nil && obj.Lon != nil {
		return *obj.Lat, *obj.Lon
	}
	if obj.Center != nil && obj.Center.Lat != 0 && obj.Center.Lon != 0 {
		return obj.Center.Lat, obj.Center.Lon
	}
	return 0, 0
}

func name(t map[string]string) *string {
	v, ok := t["name"]
	if !ok {
		return nil
	}

	return &v
}
