package job

import (
	"log/slog"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/maps"
)

const MapsName = "maps"

const sourceOverpassGate = "overpass_gate"

const datasetBusiness = "business"

type Maps struct {
	client *maps.Client
	store  *store
}

func NewMaps(
	log *slog.Logger,
	client *maps.Client,
	writer Writer,
	types TypeResolver,
	writeConcurrency int,
) *Maps {
	return &Maps{
		client: client,
		store: &store{
			log:         log.With(slog.String("source", MapsName)),
			writer:      writer,
			types:       types,
			concurrency: writeConcurrency,
		},
	}
}

func (j *Maps) InfraDatasets() []dataset {
	c := j.client
	s := j.store

	return []dataset{
		ds(datasetBusiness, sourceOverpassGate, c.BaseURL, c.GetBusinessInfra,
			s.infra(datasetBusiness)),
	}
}
