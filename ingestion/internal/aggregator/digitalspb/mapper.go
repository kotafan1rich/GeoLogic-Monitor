package digitalspb

import (
	"fmt"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const providerPrefix = "digitalspb"

const (
	datasetRailwayStation = "railway_station"
	datasetSubway         = "subway"
	datasetRestaurant     = "restaurant"
	datasetHotel          = "hotel"
	datasetMuseum         = "museum"
	datasetExhibitionHall = "exhibition_hall"
	datasetTheatre        = "theatre"
	datasetPharmacy       = "pharmacy"
	// datasetKidsPlace          = "kids_place"
	datasetLibrary         = "library"
	datasetZoo             = "zoo"
	datasetGameCenter      = "game_center"
	datasetCamp            = "camp"
	datasetEducationCenter = "education_center"
	datasetPark            = "park"
	datasetPlayground      = "playground"
	datasetSportsCenter    = "sports_center"
	datasetOther           = "other"
)

func mapToInfraObjects[T any](
	objs []T, typeID string, fn func(T) geoapi.InfraObjectInput,
) []geoapi.InfraObjectInput {
	mapped := make([]geoapi.InfraObjectInput, 0, len(objs))

	for _, obj := range objs {
		in := fn(obj)
		in.TypeID = typeID

		mapped = append(mapped, in)
	}

	return mapped
}

func mapRailwayStation(o RailwayStation) geoapi.InfraObjectInput {
	return infraObject(datasetRailwayStation, o.Number, o.Address, o.Name, o.Coordinates)
}

func mapRestaurant(o Restaurant) geoapi.InfraObjectInput {
	return infraObject(datasetRestaurant, o.OID, o.AddressManual, o.Name, o.Coord)
}

func mapHotel(o Hotel) geoapi.InfraObjectInput {
	return infraObject(datasetHotel, o.OID, o.AddressManual, o.Name, o.Coord)
}

func mapMuseum(o Museum) geoapi.InfraObjectInput {
	return infraObject(datasetMuseum, o.OID, o.AddressManual, o.Name, o.Coord)
}

func mapExhibitionHall(o ExhibitionHall) geoapi.InfraObjectInput {
	return infraObject(datasetExhibitionHall, o.OID, o.AddressManual, o.Name, o.Coord)
}

func mapTheatre(o Theatre) geoapi.InfraObjectInput {
	return infraObject(datasetTheatre, o.OID, o.AddressManual, o.Name, o.Coord)
}

func mapPharmacy(o Pharmacy) geoapi.InfraObjectInput {
	return infraObject(datasetPharmacy, o.Number, o.Address, o.Name, o.Coordinates)
}

func mapSubway(o Subway) geoapi.InfraObjectInput {
	return geoapi.InfraObjectInput{
		ExternalID: externalID(datasetSubway, o.Title),
		Address:    o.Address,
		Lat:        o.Lat,
		Lon:        o.Lon,
		Name:       optional(o.Title),
	}
}

// func mapKidsPlace(o KidsPlace) geoapi.InfraObjectInput {
// 	return infraObject()
// }

func infraObject[I int | int64 | string](
	dataset string, id I, address, name string, coord []float64,
) geoapi.InfraObjectInput {
	lat, lon := coords(coord)

	return geoapi.InfraObjectInput{
		ExternalID: externalID(dataset, id),
		Address:    address,
		Lat:        lat,
		Lon:        lon,
		Name:       optional(name),
	}
}

func externalID[I int | int64 | string](dataset string, id I) string {
	return fmt.Sprintf("%s:%s:%v", providerPrefix, dataset, id)
}

func coords(c []float64) (lat, lon float64) {
	if len(c) < 2 {
		return 0, 0
	}

	return c[0], c[1]
}

func optional(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
