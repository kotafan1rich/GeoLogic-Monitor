package digitalspb

import (
	"context"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
)

const railwayStationEndpoint = "datasets/154/versions/latest/data/188/"

func (c *Client) ParseRailwayStationData(ctx context.Context, src a.URL) ([]RailwayStation, error) {
	const op = "digitalspb.Client.ParseRailwayStationData"

	data, err := fetchSpbClassifGate[RailwayStation](ctx, c, src, railwayStationEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseSubwayData(ctx context.Context, src a.File) ([]Subway, error) {
	const op = "digitalspb.Client.ParseSubwayData"

	data, err := fetchStatic[Subway](ctx, src, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}
