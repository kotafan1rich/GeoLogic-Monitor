package geoapi

import (
	"context"
	"fmt"
	"net/url"
)

const trackedLocationsEndpoint = "internal/v1/tracked-locations"

func (c *Client) TrackedLocations(ctx context.Context) ([]TrackedLocation, error) {
	const op = "geoapi.Client.TrackedLocations"

	var raw []TrackedLocation

	if err := c.get(ctx, trackedLocationsEndpoint, url.Values{}, &raw); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}
