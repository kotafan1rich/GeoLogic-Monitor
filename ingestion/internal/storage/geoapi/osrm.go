package geoapi

import (
	"context"
	"fmt"
	"net/url"
)

const osrmEndpoint = "internal/v1/routes/walking-distances"

func (c *Client) CalculateWalkingDistance(
	ctx context.Context, src GeoPoint, dests []GeoPoint,
) ([]*float64, error) {
	const op = "geoapi.Client.CalculateWalkingDistance"

	var raw WalkingDistancesResponse

	if err := c.post(ctx, osrmEndpoint, url.Values{}, &WalkingDistancesRequest{
		Source:       src,
		Destinations: dests,
	}, &raw); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return raw.Distances, nil
}
