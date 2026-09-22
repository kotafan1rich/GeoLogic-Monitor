package geoapi

import (
	"context"
	"fmt"
	"net/url"
)

const geocodeEndpoint = "internal/v1/geocoding/suggestions"

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

	// TODO: согласовать возвращение первого элемента из слайса
	return raw[0].Lat, raw[0].Lon, nil
}
