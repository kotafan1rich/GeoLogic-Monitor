package geoapi

import (
	"context"
	"fmt"
)

const eventEndpoint = "internal/v1/events"

func (c *Client) PutEvent(ctx context.Context, obj EventInput) (Event, error) {
	const op = "geoapi.Client.PutEvent"

	var raw Event

	if err := c.put(ctx, eventEndpoint, obj, &raw); err != nil {
		return Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}
