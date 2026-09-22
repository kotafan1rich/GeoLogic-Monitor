package job

import (
	"context"

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

const (
	sourceSubwayFile = "subway"
)

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
	datasetProperty       = "property"
	datasetKidsPlace      = "kids_place"
	datasetVisit          = "visit"
	datasetStreetMusician = "street_musicians"
	datasetSubway         = "subway"
)

type DigitalSpb struct {
	client  *digitalspb.Client
	files   map[string]a.File
	writer  DigitalSpbWriter
	types   TypeResolver
	convert digitalspb.AddressConverter
	address digitalspb.CoordinatesConverter
}

func NewDigitalSpb(
	client *digitalspb.Client,
	staticFiles map[string]string,
	writer DigitalSpbWriter,
	types TypeResolver,
	convert digitalspb.AddressConverter,
	address digitalspb.CoordinatesConverter,
) *DigitalSpb {
	files := make(map[string]a.File, len(staticFiles))
	for key, path := range staticFiles {
		files[key] = a.NewFile(path)
	}

	return &DigitalSpb{
		client:  client,
		files:   files,
		writer:  writer,
		types:   types,
		convert: convert,
		address: address,
	}
}

func (j *DigitalSpb) InfraDatasets() []dataset {
	c := j.client
	w, t, conv, addr := j.writer, j.types, j.convert, j.address

	classif := j.url(sourceSpbClassifGate)
	yazzh := j.url(sourceYazzhGate)
	subway := j.file(sourceSubwayFile)

	return []dataset{
		ds(datasetRailwayStation, sourceSpbClassifGate, classif, c.ParseRailwayStationData,
			toInfra(w, t, datasetRailwayStation)),
		ds(datasetRestaurant, sourceSpbClassifGate, classif, c.ParseRestaurantData,
			toInfra(w, t, datasetRestaurant)),
		ds(datasetHotel, sourceSpbClassifGate, classif, c.ParseHotelData,
			toInfra(w, t, datasetHotel)),
		ds(datasetMuseum, sourceSpbClassifGate, classif, c.ParseMuseumData,
			toInfra(w, t, datasetMuseum)),
		ds(datasetExhibitionHall, sourceSpbClassifGate, classif, c.ParseExhibitionHallData,
			toInfra(w, t, datasetExhibitionHall)),
		ds(datasetTheatre, sourceSpbClassifGate, classif, c.ParseTheatreData,
			toInfra(w, t, datasetTheatre)),
		ds(datasetPharmacy, sourceSpbClassifGate, classif, c.ParsePharmacyData,
			toInfra(w, t, datasetPharmacy)),
		ds(datasetSubway, sourceSubwayFile, subway, c.ParseSubwayData,
			toInfra(w, t, datasetSubway)),
		ds(datasetCinema, sourceSpbClassifGate, classif, geocoded(c.ParseCinemaData, conv),
			toInfra(w, t, datasetCinema)),
		ds(datasetVetClinic, sourceSpbClassifGate, classif, geocoded(c.ParseVetClinicData, conv),
			toInfra(w, t, datasetVetClinic)),
		ds(datasetKidsPlace, sourceYazzhGate, yazzh, geocoded(c.ParseKidsPlaceData, addr),
			toInfra(w, t, datasetKidsPlace)),
	}
}

func (j *DigitalSpb) EventDatasets() []dataset {
	c := j.client
	w := j.writer

	egs := j.url(sourceEgsGate)

	return []dataset{
		ds(datasetVisit, sourceEgsGate, egs, c.ParseVisitData, toEvent(w)),
		ds(datasetStreetMusician, sourceEgsGate, egs, c.ParseStreetMusiciansData, toEvent(w)),
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
