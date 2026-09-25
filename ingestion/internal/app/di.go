package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/digitalspb"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/maps"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/closer"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/infra"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/job"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/logger"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/scheduler"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const schedulerName = "ingestion"

const (
	digitalSpbClient      = "digitalspb"
	mapsClient            = "maps-mail-ru"
	geoAPIWriteClient     = "geo-api:write"
	geoAPIGeocodingClient = "geo-api:geocoding"
)

type diContainer struct {
	log       *slog.Logger
	dsc       *digitalspb.Client
	ds        *job.DigitalSpb
	maps      *maps.Client
	m         *job.Maps
	infraJob  *job.Job
	eventsJob *job.Job
	s         *scheduler.Scheduler
	gc        *geoapi.Client
	tr        *geoapi.TypeRegistry
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Logger() *slog.Logger {
	if d.log == nil {
		cfg := config.Get()
		d.log = logger.MustNew(cfg.Logger.Level, cfg.Logger.Format, os.Stdout)

		d.log.Info("logger initialized",
			slog.String("level", cfg.Logger.Level),
			slog.String("format", cfg.Logger.Format),
		)
	}
	return d.log
}

func (d *diContainer) requester(
	name string, cfg config.HTTPClientConfig, authToken string, requireAuth bool,
) *infra.Requester {
	hc := infra.MustNewHTTPClient(infra.HTTPOptions{
		Name:                name,
		AuthToken:           authToken,
		RequireAuth:         requireAuth,
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		MaxConnsPerHost:     cfg.MaxConnsPerHost,
		RequestTimeout:      cfg.RequestTimeout,
		RateLimit: infra.RateLimit{
			RPS:   cfg.RateLimit.RPS,
			Burst: cfg.RateLimit.Burst,
		},
	})

	d.Logger().Info("http client initialized",
		slog.String("client", name),
		slog.Duration("request_timeout", cfg.RequestTimeout),
		slog.Duration("attempt_timeout", cfg.AttemptTimeout),
		slog.Uint64("max_retries", uint64(cfg.MaxRetries)),
		slog.Float64("rate_limit_rps", cfg.RateLimit.RPS),
		slog.Int("rate_limit_burst", cfg.RateLimit.Burst),
	)

	return infra.MustNewRequester(name, hc, cfg.AttemptTimeout, cfg.MaxRetries, infra.RetryTransient)
}

func (d *diContainer) DSClient() *digitalspb.Client {
	if d.dsc == nil {
		cfg := config.Get()

		req := d.requester(digitalSpbClient, cfg.Aggregator.DigitalSpb.HTTP, "", false)

		d.dsc = digitalspb.MustNew(d.Logger(), req, cfg.Aggregator.DigitalSpb.BaseURLMap)

		d.Logger().Info("digital-spb client initialized",
			slog.Int("base_urls", len(cfg.Aggregator.DigitalSpb.BaseURLMap)),
		)
	}
	return d.dsc
}

func (d *diContainer) MapsClient() *maps.Client {
	if d.maps == nil {
		cfg := config.Get()

		req := d.requester(mapsClient, cfg.Aggregator.Maps.HTTP, "", false)

		d.maps = maps.MustNew(d.Logger(), req, cfg.Aggregator.Maps.BaseURL)

		d.Logger().Info("maps-mail-ru client initialized")
	}
	return d.maps
}

func (d *diContainer) DigitalSpb() *job.DigitalSpb {
	if d.ds == nil {
		cfg := config.Get()

		d.ds = job.NewDigitalSpb(
			d.Logger(),
			d.DSClient(),
			cfg.Aggregator.DigitalSpb.StaticFiles,
			d.GeoApiClient(),
			d.TypeRegistry(),
			d.GeoApiClient().Coordinates,
			d.GeoApiClient().Address,
			cfg.GeoApi.WriteConcurrency,
		)

		d.Logger().Info(
			"datasets source initialized",
			slog.String("source", job.DigitalSpbName),
			slog.Int("static_files", len(cfg.Aggregator.DigitalSpb.StaticFiles)),
			slog.Int("write_concurrency", cfg.GeoApi.WriteConcurrency),
		)
	}
	return d.ds
}

func (d *diContainer) Maps() *job.Maps {
	if d.m == nil {
		cfg := config.Get()

		d.m = job.NewMaps(
			d.Logger(),
			d.MapsClient(),
			d.GeoApiClient(),
			d.TypeRegistry(),
			cfg.GeoApi.WriteConcurrency,
		)

		d.Logger().Info(
			"datasets source initialized",
			slog.String("source", job.MapsName),
			slog.Int("write_concurrency", cfg.GeoApi.WriteConcurrency),
		)
	}
	return d.m
}

func (d *diContainer) InfraJob() *job.Job {
	if d.infraJob == nil {
		cfg := config.Get()

		d.infraJob = job.New(
			job.InfraName,
			cfg.Scheduler.Infra,
			d.Logger(),
			d.DigitalSpb().InfraDatasets,
			d.Maps().InfraDatasets,
		)

		d.logJob(d.infraJob)
	}
	return d.infraJob
}

// TODO: Впоследствии дополнить cudago (Опционально)
func (d *diContainer) EventsJob() *job.Job {
	if d.eventsJob == nil {
		cfg := config.Get()

		d.eventsJob = job.New(
			job.EventsName,
			cfg.Scheduler.Events,
			d.Logger(),
			d.DigitalSpb().EventDatasets,
		)

		d.logJob(d.eventsJob)
	}
	return d.eventsJob
}

func (d *diContainer) Jobs() []scheduler.Job {
	return []scheduler.Job{d.InfraJob(), d.EventsJob()}
}

func (d *diContainer) logJob(j *job.Job) {
	d.Logger().Info(
		"job initialized",
		slog.String("job", j.Name()),
		slog.String("schedule", j.Schedule()),
	)
}

func (d *diContainer) Scheduler() *scheduler.Scheduler {
	if d.s == nil {
		s := scheduler.MustNew(schedulerName, d.Logger())

		d.Logger().Info(
			"scheduler initialized",
			slog.String("scheduler", s.Name()),
		)

		closer.Add("scheduler:"+s.Name(), func(context.Context) error {
			return s.Shutdown()
		})

		d.s = s
	}
	return d.s
}

func (d *diContainer) GeoApiClient() *geoapi.Client {
	if d.gc == nil {
		cfg := config.Get()

		write := d.requester(geoAPIWriteClient, cfg.GeoApi.Write, cfg.GeoApi.AuthToken, true)
		geocoding := d.requester(geoAPIGeocodingClient, cfg.GeoApi.Geocoding, cfg.GeoApi.AuthToken, true)

		d.gc = geoapi.MustNew(cfg.GeoApi.URL, write, geocoding)

		d.Logger().Info("geo-api client initialized", slog.String("url", cfg.GeoApi.URL))
	}
	return d.gc
}

func (d *diContainer) TypeRegistry() *geoapi.TypeRegistry {
	if d.tr == nil {
		cfg := config.Get()

		defs := make([]geoapi.InfraTypeInput, 0, len(cfg.GeoApi.InfraTypes))
		for _, t := range cfg.GeoApi.InfraTypes {
			defs = append(defs, geoapi.InfraTypeInput{
				Slug:      t.Slug,
				Name:      t.Name,
				Weight:    t.Weight,
				MaxRadius: t.MaxRadius,
			})
		}

		d.tr = geoapi.MustNewTypeRegistry(d.GeoApiClient(), defs)

		d.Logger().Info("infra type registry initialized", slog.Int("types", len(defs)))
	}
	return d.tr
}
