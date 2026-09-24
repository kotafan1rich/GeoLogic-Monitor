package digitalspb

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const providerPrefix = "digitalspb"

const geocodeLogEvery = 250

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
	datasetKidsPlace       = "kids_place"
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

	sourceEgsGate = "egs_gate"

	datasetVisit          = "visit"
	datasetStreetMusician = "street_musicians"

	coordSep = ","
)

func mapToInfraObjects[T any](objs []T, fn func(T) geoapi.InfraObjectInput) []geoapi.InfraObjectInput {
	mapped := make([]geoapi.InfraObjectInput, 0, len(objs))

	for _, obj := range objs {
		mapped = append(mapped, fn(obj))
	}

	return mapped
}

type AddressConverter func(ctx context.Context, address string) (lat, lon float64, err error)

type CoordinatesConverter func(ctx context.Context, lat, lon float64) (address string, err error)

func mapToGeocodedInfraObjects[T any](
	ctx context.Context,
	log *slog.Logger,
	dataset string,
	objs []T,
	convert AddressConverter,
	fn func(T) geoapi.InfraObjectInput,
) ([]geoapi.InfraObjectInput, error) {
	if convert == nil {
		return nil, ErrInvalidAddressConverter
	}

	mapped := mapToInfraObjects(objs, fn)

	err := geocode(
		ctx, log, dataset, "address_to_coordinates", mapped,
		func(o geoapi.InfraObjectInput) bool {
			return o.Address != "" && o.Lat == 0 && o.Lon == 0
		},
		func(ctx context.Context, o *geoapi.InfraObjectInput) error {
			lat, lon, err := convert(ctx, o.Address)
			if err != nil {
				return err
			}

			o.Lat, o.Lon = lat, lon

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return mapped, nil
}

func mapToAddressedInfraObjects[T any](
	ctx context.Context,
	log *slog.Logger,
	dataset string,
	objs []T,
	convert CoordinatesConverter,
	fn func(T) geoapi.InfraObjectInput,
) ([]geoapi.InfraObjectInput, error) {
	if convert == nil {
		return nil, ErrInvalidCoordinatesConverter
	}

	mapped := mapToInfraObjects(objs, fn)

	err := geocode(
		ctx, log, dataset, "coordinates_to_address", mapped,
		func(o geoapi.InfraObjectInput) bool {
			return o.Address == "" && (o.Lat != 0 || o.Lon != 0)
		},
		func(ctx context.Context, o *geoapi.InfraObjectInput) error {
			address, err := convert(ctx, o.Lat, o.Lon)
			if err != nil {
				return err
			}

			o.Address = address

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return mapped, nil
}

func geocode(
	ctx context.Context,
	log *slog.Logger,
	dataset, direction string,
	objs []geoapi.InfraObjectInput,
	needs func(geoapi.InfraObjectInput) bool,
	fill func(context.Context, *geoapi.InfraObjectInput) error,
) error {
	total := 0

	for i := range objs {
		if needs(objs[i]) {
			total++
		}
	}

	if total == 0 {
		return nil
	}

	log = log.With(
		slog.String("dataset", dataset),
		slog.String("direction", direction),
	)

	log.DebugContext(ctx, "geocoding started", slog.Int("objects", total))

	var (
		start     = time.Now()
		collector = errs.NewCollector(errs.DefaultSampleSize)
		processed int
		resolved  int
	)

	for i := range objs {
		if !needs(objs[i]) {
			continue
		}

		if err := ctx.Err(); err != nil {
			log.WarnContext(ctx, "geocoding interrupted",
				slog.Int("processed", processed),
				slog.Int("objects", total),
			)

			return err
		}

		processed++

		if err := fill(ctx, &objs[i]); err != nil {
			collector.Add(err)
			continue
		}

		resolved++

		if processed%geocodeLogEvery == 0 {
			log.DebugContext(ctx, "geocoding in progress",
				slog.Int("processed", processed),
				slog.Int("objects", total),
				slog.Int("failed", collector.Total()),
			)
		}
	}

	attrs := []any{
		slog.Int("objects", total),
		slog.Int("resolved", resolved),
		slog.Int("failed", collector.Total()),
		slog.Duration("duration", time.Since(start)),
	}

	if collector.Total() > 0 {
		log.WarnContext(ctx, "geocoding finished with errors", append(attrs, collector.Attr())...)

		return nil
	}

	log.DebugContext(ctx, "geocoding finished", attrs...)

	return nil
}

func mapToEvents[T any](objs []T, fn func(T) geoapi.EventInput) []geoapi.EventInput {
	mapped := make([]geoapi.EventInput, 0, len(objs))

	for _, obj := range objs {
		mapped = append(mapped, fn(obj))
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

var kidsPlaceDatasets = map[string]string{
	kidsPlaceLibrary:         datasetLibrary,
	kidsPlaceZoo:             datasetZoo,
	kidsPlaceGameCenter:      datasetGameCenter,
	kidsPlaceCamp:            datasetCamp,
	kidsPlaceMuseum:          datasetMuseum,
	kidsPlaceEducationCenter: datasetEducationCenter,
	kidsPlacePark:            datasetPark,
	kidsPlacePlayground:      datasetPlayground,
	kidsPlaceSportsCenter:    datasetSportsCenter,
	kidsPlaceTheatre:         datasetTheatre,
}

func mapKidsPlace(o KidsPlace) geoapi.InfraObjectInput {
	dataset, ok := kidsPlaceDatasets[o.CategoriesName]
	if !ok {
		dataset = datasetOther
	}

	return infraObject(dataset, o.ID, "", o.Title, o.Coordinates)
}

func mapStreetPerformance(o StreetPerformance) geoapi.EventInput {
	return event(datasetStreetMusician, sourceEgsGate, o.Address, o.ID, o.Coordinates, extractDate(o.StartDate))
}

func mapCultureEvent(o CultureEvent) geoapi.EventInput {
	return event(datasetVisit, sourceEgsGate, o.Name, o.ID, extractCoords(o.Map), extractDate(o.Start))
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

func event[I int | int64 | string](
	dataset, provider, info string, id I, coord []float64, date time.Time,
) geoapi.EventInput {
	lat, lon := coords(coord)

	return geoapi.EventInput{
		Provider:   provider,
		ExternalID: externalID(dataset, id),
		Lat:        lat,
		Lon:        lon,
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
		return nil
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(raw[0]), 64)
	if err != nil {
		return nil
	}

	lon, err := strconv.ParseFloat(strings.TrimSpace(raw[1]), 64)
	if err != nil {
		return nil
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
