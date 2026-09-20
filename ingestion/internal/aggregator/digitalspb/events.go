package digitalspb

import (
	"context"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
)

const (
	visitEndpoint           = "visti_spb/"
	streetMusiciansEndpoint = "street_musicians/external/event/"
)

func (c *Client) ParseVisitData(ctx context.Context, src a.URL) ([]CultureEvent, error) {
	const op = "digitalspb.Client.ParseVisitData"

	data, err := fetchEgsGate[CultureEvent](ctx, c, src, egsGateV1, visitEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseStreetMusiciansData(ctx context.Context, src a.URL) ([]StreetPerformance, error) {
	const op = "digitalspb.Client.ParseStreetMusiciansData"

	data, err := fetchEgsGate[StreetPerformance](ctx, c, src, egsGateV2, streetMusiciansEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}
