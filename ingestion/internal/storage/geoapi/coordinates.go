package geoapi

import (
	"context"
	"fmt"
	"net/url"
)

const (
	geocodeEndpoint        = "internal/v1/geocoding/suggestions"
	reverseGeocodeEndpoint = "internal/v1/geocoding/address"
)

func (c *Client) Coordinates(ctx context.Context, address string) (float64, float64, error) {
	const op = "geoapi.Client.Coordinates"

	query := url.Values{"query": {address}}

	var raw []AddressComponent

	if err := c.get(ctx, geocodeEndpoint, query, &raw); err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	if len(raw) == 0 {
		return 0, 0, fmt.Errorf("%s: %w [%s]", op, ErrEmptyAddressList, address)
	}

	return raw[0].Lat, raw[0].Lon, nil
}
