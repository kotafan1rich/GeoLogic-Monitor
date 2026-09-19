package job

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/digitalspb"
)

const DigitalSpbName = "digitalspb"

const (
	sourceSpbClassifGate = "spb_classif_gate"
	sourceEgsGate        = "egs_gate"
	sourceYazzhGate      = "yazzh_gate"
)

type DigitalSpb struct {
	log      *slog.Logger
	client   *digitalspb.Client
	interval time.Duration
}

func NewDigitalSpb(log *slog.Logger, client *digitalspb.Client, interval time.Duration) *DigitalSpb {
	return &DigitalSpb{log: log, client: client, interval: interval}
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

	return []dataset{
		j.ds("railway_station", sourceSpbClassifGate, collect(c.ParseRailwayStationData)),
		j.ds("restaurant", sourceSpbClassifGate, collect(c.ParseRestaurantData)),
		j.ds("hotel", sourceSpbClassifGate, collect(c.ParseHotelData)),
		j.ds("cinema", sourceSpbClassifGate, collect(c.ParseCinemaData)),
		j.ds("museum", sourceSpbClassifGate, collect(c.ParseMuseumData)),
		j.ds("exhibition_hall", sourceSpbClassifGate, collect(c.ParseExhibitionHallData)),
		j.ds("theatre", sourceSpbClassifGate, collect(c.ParseTheatreData)),
		j.ds("pharmacy", sourceSpbClassifGate, collect(c.ParsePharmacyData)),
		j.ds("vet_clinic", sourceSpbClassifGate, collect(c.ParseVetClinicData)),
		j.ds("property", sourceSpbClassifGate, collect(c.ParsePropertyData)),
		j.ds("kids_place", sourceYazzhGate, collect(c.ParseKidsPlace)),
		j.ds("visit", sourceEgsGate, collect(c.ParseVisitData)),
		j.ds("street_musicians", sourceEgsGate, collect(c.ParseStreetMusiciansData)),
	}
}

func (j *DigitalSpb) ds(name, source string, parse parseFunc) dataset {
	var baseURL *url.URL
	if u, ok := j.client.BaseURLs[source]; ok {
		baseURL = u
	}

	return dataset{name: name, source: source, baseURL: baseURL, parse: parse}
}
