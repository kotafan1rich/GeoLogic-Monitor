package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/digitalspb"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/closer"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/infra"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/job"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/logger"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/scheduler"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const schedulerName = "ingestion"

type diContainer struct {
	log       *slog.Logger
	dsc       *digitalspb.Client
	ds        *job.DigitalSpb
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

func (d *diContainer) DSClient() *digitalspb.Client {
	if d.dsc == nil {
		cfg := config.Get()
		log := d.Logger()

		httpClient := infra.NewDigitalSpbHTTPClient(
			cfg.Aggregator.DigitalSpb.MaxIdleConns,
			cfg.Aggregator.DigitalSpb.MaxIdleConnsPerHost,
			cfg.Aggregator.DigitalSpb.MaxConnsPerHost,
			cfg.Aggregator.DigitalSpb.RequestTimeout,
		)

		log.Info("digital-spb http client initialized")

		d.dsc = digitalspb.MustNew(
			httpClient,
			cfg.Aggregator.DigitalSpb.BaseURLMap,
			cfg.Aggregator.DigitalSpb.AttemptTimeout,
			cfg.Aggregator.DigitalSpb.MaxRetries,
		)

		log.Info("digital-spb client initialized")
	}
	return d.dsc
}

func (d *diContainer) DigitalSpb() *job.DigitalSpb {
	if d.ds == nil {
		cfg := config.Get()

		d.ds = job.NewDigitalSpb(
			d.DSClient(),
			cfg.Aggregator.DigitalSpb.StaticFiles,
			d.GeoApiClient(),
			d.TypeRegistry(),
			d.GeoApiClient().Coordinates,
		)

		d.Logger().Info(
			"datasets source initialized",
			slog.String("source", job.DigitalSpbName),
			slog.Int("static_files", len(cfg.Aggregator.DigitalSpb.StaticFiles)),
		)
	}
	return d.ds
}

func (d *diContainer) InfraJob() *job.Job {
	if d.infraJob == nil {
		cfg := config.Get()

		d.infraJob = job.New(
			job.InfraName,
			cfg.Scheduler.Infra,
			d.Logger(),
			d.DigitalSpb().InfraDatasets,
		)

		d.logJob(d.infraJob)
	}
	return d.infraJob
}

// TODO: Впоследствии добавить cudago (Опционально)
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
		log := d.Logger()

		httpClient := infra.MustNewGeoAPIHTTPClient(
			cfg.GeoApi.AuthToken,
			cfg.GeoApi.MaxIdleConns,
			cfg.GeoApi.MaxIdleConnsPerHost,
			cfg.GeoApi.MaxConnsPerHost,
			cfg.GeoApi.RequestTimeout,
		)

		log.Info("geo-api http client initialized")

		d.gc = geoapi.MustNew(
			httpClient,
			cfg.GeoApi.URL,
			cfg.GeoApi.AttemptTimeout,
			cfg.GeoApi.MaxRetries,
		)

		log.Info("geo-api client initialized")
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
