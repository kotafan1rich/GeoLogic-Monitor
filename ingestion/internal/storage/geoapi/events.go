package geoapi

import (
	"context"
	"fmt"
	"net/url"
)

const eventEndpoint = "internal/v1/events"

func (c *Client) PutEvent(ctx context.Context, obj EventInput) (Event, error) {
	const op = "geoapi.Client.PutEvent"

	var raw Event

	err := c.do(ctx, c.baseURL, eventEndpoint, url.Values{}, obj, &raw)
	if err != nil {
		return Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}
