package digitalspb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const providerPrefix = "digitalspb"

const geocodeConcurrency = 8

const (
	datasetRailwayStation  = "railway_station"
	datasetSubway          = "subway"
	datasetRestaurant      = "restaurant"
	datasetHotel           = "hotel"
	datasetMuseum          = "museum"
	datasetExhibitionHall  = "exhibition_hall"
	datasetTheatre         = "theatre"
	datasetPharmacy        = "pharmacy"
	datasetCinema          = "cinema"
	datasetVetClinic       = "vet_clinic"
	datasetLibrary         = "library"
	datasetZoo             = "zoo"
	datasetGameCenter      = "game_center"
	datasetCamp            = "camp"
	datasetEducationCenter = "education_center"
	datasetPark            = "park"
	datasetPlayground      = "playground"
	datasetSportsCenter    = "sports_center"
	datasetOther           = "other"

	kidsPlaceLibrary         = "Библиотеки"
	kidsPlaceZoo             = "Зоопарки"
	kidsPlaceGameCenter      = "Игровые центры"
	kidsPlaceCamp            = "Лагеря"
	kidsPlaceMuseum          = "Музеи"
	kidsPlaceEducationCenter = "Образовательные центры"
	kidsPlacePark            = "Парки"
	kidsPlacePlayground      = "Площадки"
	kidsPlaceSportsCenter    = "Спортивные центры"
	kidsPlaceTheatre         = "Театры"
	kidsPlaceOther           = "Другое"

	sourceEgsGate = "egs_gate"

	datasetVisit          = "visit"
	datasetStreetMusician = "street_musicians"

	coordSep = ","
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

type AddressConverter func(ctx context.Context, address string) (lat, lon float64, err error)

func mapToGeocodedInfraObjects[T any](
	ctx context.Context, objs []T, typeID string, convert AddressConverter, fn func(T) geoapi.InfraObjectInput,
) ([]geoapi.InfraObjectInput, error) {
	if convert == nil {
		return nil, ErrInvalidAddressConverter
	}

	mapped := mapToInfraObjects(objs, typeID, fn)

	var (
		g      errgroup.Group
		failed atomic.Int64
		mu     sync.Mutex
		errs   []error
	)

	g.SetLimit(geocodeConcurrency)

	for i := range mapped {
		if ctx.Err() != nil {
			break
		}

		if mapped[i].Address == "" {
			continue
		}

		idx := i

		g.Go(func() error {
			if err := ctx.Err(); err != nil {
				failed.Add(1)
				return nil
			}

			lat, lon, err := convert(ctx, mapped[idx].Address)
			if err != nil {
				failed.Add(1)
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return nil
			}

			mapped[idx].Lat, mapped[idx].Lon = lat, lon

			return nil
		})
	}

	_ = g.Wait()

	if failed.Load() > 0 && len(errs) > 0 {
		slog.Error(
			"addresses converting finished with errors",
			slog.Int64("failed", failed.Load()),
			slog.Any("error", errors.Join(errs...)),
		)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return mapped, nil
}

func mapToEvents[T any](objs []T, fn func(T) geoapi.EventInput) []geoapi.EventInput {
	mapped := make([]geoapi.EventInput, 0, len(objs))

	for _, obj := range objs {
		in := fn(obj)
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

func mapCinema(o Cinema) geoapi.InfraObjectInput {
	return infraObject(datasetCinema, o.Number, o.Address, o.Name, nil)
}

func mapVetClinic(o VetClinic) geoapi.InfraObjectInput {
	return infraObject(datasetVetClinic, o.Number, o.Address, o.Name, nil)
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

// TODO: потом нужно конвертировать координаты в строковый адрес
func mapKidsPlaceLibrary(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetLibrary, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlaceZoo(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetZoo, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlaceGameCenter(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetGameCenter, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlaceCamp(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetCamp, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlaceMuseum(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetMuseum, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlaceEducationCenter(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetEducationCenter, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlacePark(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetPark, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlacePlayground(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetPlayground, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlaceSportsCenter(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetSportsCenter, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlacekidsPlaceTheatre(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetTheatre, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlacekidsOther(o KidsPlace) geoapi.InfraObjectInput {
	return infraObject(datasetOther, o.ID, "test", o.Title, o.Coordinates)
}

func mapKidsPlace(o KidsPlace) geoapi.InfraObjectInput {
	switch o.CategoriesName {
	case kidsPlaceLibrary:
		return mapKidsPlaceLibrary(o)
	case kidsPlaceZoo:
		return mapKidsPlaceZoo(o)
	case kidsPlaceGameCenter:
		return mapKidsPlaceGameCenter(o)
	case kidsPlaceCamp:
		return mapKidsPlaceCamp(o)
	case kidsPlaceMuseum:
		return mapKidsPlaceMuseum(o)
	case kidsPlaceEducationCenter:
		return mapKidsPlaceEducationCenter(o)
	case kidsPlacePark:
		return mapKidsPlacePark(o)
	case kidsPlacePlayground:
		return mapKidsPlacePlayground(o)
	case kidsPlaceSportsCenter:
		return mapKidsPlaceSportsCenter(o)
	case kidsPlaceTheatre:
		return mapKidsPlacekidsPlaceTheatre(o)
	default:
		return mapKidsPlacekidsOther(o)
	}
}

func mapStreetPerformance(o StreetPerformance) geoapi.EventInput {
	return Event(
		datasetStreetMusician, sourceEgsGate, o.Address, o.ID, o.Coordinates, extractDate(o.StartDate),
	)
}

func mapCultureEvent(o CultureEvent) geoapi.EventInput {
	return Event(
		datasetVisit, sourceEgsGate, o.Name, o.ID, extractCoords(o.Map), extractDate(o.Start),
	)
}

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

func Event[I int | int64 | string, C string | []float64](
	dataset, provider, info string, id I, coord C, date time.Time,
) geoapi.EventInput {
	var extractedCoord []float64

	switch v := any(coord).(type) {
	case string:
		extractedCoord = extractCoords(v)
	case []float64:
		extractedCoord = v
	}

	return geoapi.EventInput{
		Provider:   provider,
		ExternalID: externalID(dataset, id),
		Lat:        extractedCoord[0],
		Lon:        extractedCoord[1],
		Date:       date,
		Info:       optional(info),
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

func extractCoords(c string) []float64 {
	raw := strings.Split(c, coordSep)
	if len(raw) < 2 {
		return []float64{0, 0}
	}

	lat, err := strconv.ParseFloat(raw[0], 64)
	if err != nil {
		return []float64{0, 0}
	}

	lon, err := strconv.ParseFloat(raw[0], 64)
	if err != nil {
		return []float64{0, 0}
	}

	return []float64{lat, lon}
}

func extractDate(d string) time.Time {
	layouts := []string{
		"02.01.2006 15:04:05", time.RFC3339, time.RFC3339Nano,
		"2006-01-02 15:04:05", "2006-01-02", "02.01.2006",
	}

	for _, l := range layouts {
		if date, err := time.Parse(l, d); err == nil {
			return date
		}
	}

	return time.Time{}
}

func optional(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
