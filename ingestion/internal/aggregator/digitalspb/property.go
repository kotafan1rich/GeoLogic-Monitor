package digitalspb

import (
	"context"
	"net/url"
)

const propertyEndpoint = "datasets/208/versions/latest/data/238/"

func (c *Client) ParsePropertyData(ctx context.Context, baseURL *url.URL) ([]MSPProperty, error) {
	const op = "digitalspb.Client.ParsePropertyData"

	data, err := fetchSpbClassifGate[MSPProperty](ctx, c, baseURL, propertyEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}
