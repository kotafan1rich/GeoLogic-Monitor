package digitalspb

import (
	"context"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
)

const propertyEndpoint = "datasets/208/versions/latest/data/238/"

func (c *Client) ParsePropertyData(ctx context.Context, src a.URL) ([]MSPProperty, error) {
	const op = "digitalspb.Client.ParsePropertyData"

	data, err := fetchSpbClassifGate[MSPProperty](ctx, c, src, propertyEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}
