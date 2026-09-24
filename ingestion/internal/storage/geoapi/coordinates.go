package geoapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	geocodeEndpoint        = "internal/v1/geocoding/suggestions"
	reverseGeocodeEndpoint = "internal/v1/geocoding/address"
)

func (c *Client) Coordinates(ctx context.Context, address string) (float64, float64, error) {
	const op = "geoapi.Client.Coordinates"

	query := url.Values{
		"query": {address},
	}

	var raw []AddressComponent

	if err := c.do(ctx, c.baseURL, geocodeEndpoint, query, nil, &raw); err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	if len(raw) == 0 {
		return 0, 0, fmt.Errorf("%s: %w [%s]", op, ErrEmptyAddressList, address)
	}

	return raw[0].Lat, raw[0].Lon, nil
}

func (c *Client) Address(ctx context.Context, lat, lon float64) (string, error) {
	const op = "geoapi.Client.Address"

	query := url.Values{
		"lat": {strconv.FormatFloat(lat, 'f', -1, 64)},
		"lon": {strconv.FormatFloat(lon, 'f', -1, 64)},
	}

	var raw AddressComponent

	if err := c.do(ctx, c.baseURL, reverseGeocodeEndpoint, query, nil, &raw); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if raw.Address == "" {
		return "", fmt.Errorf("%s: %w [%f, %f]", op, ErrEmptyAddress, lat, lon)
	}

	return raw.Address, nil
}
