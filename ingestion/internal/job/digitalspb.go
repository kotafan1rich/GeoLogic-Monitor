package job

import (
	"context"
	"log/slog"
	"time"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/digitalspb"
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
	datasetProperty       = "property"
	datasetKidsPlace      = "kids_place"
	datasetVisit          = "visit"
	datasetStreetMusician = "street_musicians"
	datasetSubway         = "subway"
)

type DigitalSpb struct {
	log      *slog.Logger
	client   *digitalspb.Client
	interval time.Duration
	files    map[string]a.File
	writer   InfraWriter
	types    TypeResolver
}

func NewDigitalSpb(
	log *slog.Logger,
	client *digitalspb.Client,
	interval time.Duration,
	staticFiles map[string]string,
	writer InfraWriter,
	types TypeResolver,
) *DigitalSpb {
	files := make(map[string]a.File, len(staticFiles))
	for key, path := range staticFiles {
		files[key] = a.NewFile(path)
	}

	return &DigitalSpb{
		log:      log,
		client:   client,
		interval: interval,
		files:    files,
		writer:   writer,
		types:    types,
	}
}

func (j *DigitalSpb) Name() string {
	return DigitalSpbName
}

func (j *DigitalSpb) Interval() time.Duration {
	return j.interval
}

func (j *DigitalSpb) Run(ctx context.Context) {
	runDatasets(ctx, j.log, j.Name(), j.datasets())
}

func (j *DigitalSpb) datasets() []dataset {
	c := j.client
	w, t := j.writer, j.types

	classif := j.url(sourceSpbClassifGate)
	// egs := j.url(sourceEgsGate)
	// yazzh := j.url(sourceYazzhGate)
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

		// TODO: тип зависит от категории объекта, а не от датасета
		// ds(datasetKidsPlace, sourceYazzhGate, yazzh, c.ParseKidsPlaceData, discard),

		// TODO: нет координат, нужен геокодинг до записи
		// ds(datasetProperty, sourceSpbClassifGate, classif, c.ParsePropertyData, discard),

		// TODO: запись в /internal/v1/events
		// ds(datasetVisit, sourceEgsGate, egs, c.ParseVisitData, discard),
		// ds(datasetStreetMusician, sourceEgsGate, egs, c.ParseStreetMusiciansData, discard),
	}
}

func (j *DigitalSpb) url(key string) a.URL {
	return j.client.BaseURLs[key]
}

func (j *DigitalSpb) file(key string) a.File {
	return j.files[key]
}
