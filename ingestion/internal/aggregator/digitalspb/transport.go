package digitalspb

import (
	"context"
	"net/url"
)

const railwayStationEndpoint = "datasets/154/versions/latest/data/188/"

func (c *Client) ParseRailwayStationData(ctx context.Context, baseURL *url.URL) ([]RailwayStation, error) {
	const op = "digitalspb.Client.ParseRailwayStationData"

	data, err := fetchSpbClassifGate[RailwayStation](ctx, c, baseURL, railwayStationEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}
