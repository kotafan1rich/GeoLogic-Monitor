package digitalspb

import (
	"context"
	"net/url"
)

const (
	visitEndpoint           = "visti_spb/"
	streetMusiciansEndpoint = "street_musicians/external/event/"
)

func (c *Client) ParseVisitData(ctx context.Context, baseURL *url.URL) ([]CultureEvent, error) {
	const op = "digitalspb.Client.ParseVisitData"

	data, err := fetchEgsGate[CultureEvent](ctx, c, baseURL, egsGateV1, visitEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseStreetMusiciansData(ctx context.Context, baseURL *url.URL) ([]StreetPerformance, error) {
	const op = "digitalspb.Client.ParseStreetMusiciansData"

	data, err := fetchEgsGate[StreetPerformance](ctx, c, baseURL, egsGateV2, streetMusiciansEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}
