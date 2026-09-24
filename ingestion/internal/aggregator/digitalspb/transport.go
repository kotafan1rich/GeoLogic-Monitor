package digitalspb

import (
	"context"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const railwayStationEndpoint = "datasets/154/versions/latest/data/188/"

func (c *Client) ParseRailwayStationData(ctx context.Context, src a.URL) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseRailwayStationData"

	data, err := fetchSpbClassifGate[RailwayStation](ctx, c, src, railwayStationEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData := mapToInfraObjects[RailwayStation](data, "", mapRailwayStation)

	return mappedData, nil
}

func (c *Client) ParseSubwayData(ctx context.Context, src a.File) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseSubwayData"

	data, err := fetchStatic[Subway](ctx, src, op)
	if err != nil {
		return nil, err
	}

	mappedData := mapToInfraObjects[Subway](data, "", mapSubway)

	return mappedData, nil
}
