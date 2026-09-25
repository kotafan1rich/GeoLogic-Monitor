package job

import (
	"context"
	"log/slog"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/digitalspb"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const DigitalSpbName = "digitalspb"

const (
	sourceSpbClassifGate = "spb_classif_gate"
	sourceEgsGate        = "egs_gate"
	sourceYazzhGate      = "yazzh_gate"
)

const sourceSubwayFile = "subway"

const (
	datasetRailwayStation = "railway_station"
	datasetRestaurant     = "restaurant"
	datasetHotel          = "hotel"
	datasetMuseum         = "museum"
	datasetExhibitionHall = "exhibition_hall"
	datasetTheatre        = "theatre"
	datasetPharmacy       = "pharmacy"
	datasetCinema         = "cinema"
	datasetVetClinic      = "vet_clinic"
	datasetKidsPlace      = "kids_place"
	datasetVisit          = "visit"
	datasetStreetMusician = "street_musicians"
	datasetSubway         = "subway"
)

type DigitalSpb struct {
	client  *digitalspb.Client
	files   map[string]a.File
	store   *store
	convert digitalspb.AddressConverter
	address digitalspb.CoordinatesConverter
}

func NewDigitalSpb(
	log *slog.Logger,
	client *digitalspb.Client,
	staticFiles map[string]string,
	writer Writer,
	types TypeResolver,
	convert digitalspb.AddressConverter,
	address digitalspb.CoordinatesConverter,
	writeConcurrency int,
) *DigitalSpb {
	files := make(map[string]a.File, len(staticFiles))
	for key, path := range staticFiles {
		files[key] = a.NewFile(path)
	}

	return &DigitalSpb{
		client: client,
		files:  files,
		store: &store{
			log:         log.With(slog.String("source", DigitalSpbName)),
			writer:      writer,
			types:       types,
			concurrency: writeConcurrency,
		},
		convert: convert,
		address: address,
	}
}

func (j *DigitalSpb) InfraDatasets() []dataset {
	c := j.client
	s, conv, addr := j.store, j.convert, j.address

	classif := j.url(sourceSpbClassifGate)
	yazzh := j.url(sourceYazzhGate)
	subway := j.file(sourceSubwayFile)

	return []dataset{
		ds(datasetRailwayStation, sourceSpbClassifGate, classif, c.ParseRailwayStationData,
			s.infra(datasetRailwayStation)),
		ds(datasetRestaurant, sourceSpbClassifGate, classif, c.ParseRestaurantData,
			s.infra(datasetRestaurant)),
		ds(datasetHotel, sourceSpbClassifGate, classif, c.ParseHotelData,
			s.infra(datasetHotel)),
		ds(datasetMuseum, sourceSpbClassifGate, classif, c.ParseMuseumData,
			s.infra(datasetMuseum)),
		ds(datasetExhibitionHall, sourceSpbClassifGate, classif, c.ParseExhibitionHallData,
			s.infra(datasetExhibitionHall)),
		ds(datasetTheatre, sourceSpbClassifGate, classif, c.ParseTheatreData,
			s.infra(datasetTheatre)),
		ds(datasetPharmacy, sourceSpbClassifGate, classif, c.ParsePharmacyData,
			s.infra(datasetPharmacy)),
		ds(datasetSubway, sourceSubwayFile, subway, c.ParseSubwayData,
			s.infra(datasetSubway)),
		ds(datasetCinema, sourceSpbClassifGate, classif, geocoded(c.ParseCinemaData, conv),
			s.infra(datasetCinema)),
		ds(datasetVetClinic, sourceSpbClassifGate, classif, geocoded(c.ParseVetClinicData, conv),
			s.infra(datasetVetClinic)),
		ds(datasetKidsPlace, sourceYazzhGate, yazzh, geocoded(c.ParseKidsPlaceData, addr),
			s.infra(datasetKidsPlace)),
	}
}

func (j *DigitalSpb) EventDatasets() []dataset {
	c := j.client
	s := j.store

	egs := j.url(sourceEgsGate)

	return []dataset{
		ds(datasetVisit, sourceEgsGate, egs, c.ParseVisitData, s.events(datasetVisit)),
		ds(datasetStreetMusician, sourceEgsGate, egs, c.ParseStreetMusiciansData,
			s.events(datasetStreetMusician)),
	}
}

func geocoded[S a.Source, C digitalspb.AddressConverter | digitalspb.CoordinatesConverter](
	fetch func(context.Context, S, C) ([]geoapi.InfraObjectInput, error),
	convert C,
) func(context.Context, S) ([]geoapi.InfraObjectInput, error) {
	return func(ctx context.Context, src S) ([]geoapi.InfraObjectInput, error) {
		return fetch(ctx, src, convert)
	}
}

func (j *DigitalSpb) url(key string) a.URL {
	return j.client.BaseURLs[key]
}

func (j *DigitalSpb) file(key string) a.File {
	return j.files[key]
}
