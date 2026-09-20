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

type DigitalSpb struct {
	log      *slog.Logger
	client   *digitalspb.Client
	interval time.Duration
	files    map[string]a.File
}

func NewDigitalSpb(
	log *slog.Logger, client *digitalspb.Client, interval time.Duration, staticFiles map[string]string,
) *DigitalSpb {
	files := make(map[string]a.File, len(staticFiles))
	for key, path := range staticFiles {
		files[key] = a.NewFile(path)
	}

	return &DigitalSpb{log: log, client: client, interval: interval, files: files}
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

	classif := j.url(sourceSpbClassifGate)
	egs := j.url(sourceEgsGate)
	yazzh := j.url(sourceYazzhGate)
	subway := j.file(sourceSubwayFile)

	return []dataset{
		ds("railway_station", sourceSpbClassifGate, classif, c.ParseRailwayStationData),
		ds("restaurant", sourceSpbClassifGate, classif, c.ParseRestaurantData),
		ds("hotel", sourceSpbClassifGate, classif, c.ParseHotelData),
		ds("cinema", sourceSpbClassifGate, classif, c.ParseCinemaData),
		ds("museum", sourceSpbClassifGate, classif, c.ParseMuseumData),
		ds("exhibition_hall", sourceSpbClassifGate, classif, c.ParseExhibitionHallData),
		ds("theatre", sourceSpbClassifGate, classif, c.ParseTheatreData),
		ds("pharmacy", sourceSpbClassifGate, classif, c.ParsePharmacyData),
		ds("vet_clinic", sourceSpbClassifGate, classif, c.ParseVetClinicData),
		ds("property", sourceSpbClassifGate, classif, c.ParsePropertyData),
		ds("kids_place", sourceYazzhGate, yazzh, c.ParseKidsPlaceData),
		ds("visit", sourceEgsGate, egs, c.ParseVisitData),
		ds("street_musicians", sourceEgsGate, egs, c.ParseStreetMusiciansData),
		ds("subway", sourceSubwayFile, subway, c.ParseSubwayData),
	}
}

func (j *DigitalSpb) url(key string) a.URL {
	return j.client.BaseURLs[key]
}

func (j *DigitalSpb) file(key string) a.File {
	return j.files[key]
}
