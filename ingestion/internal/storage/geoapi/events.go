package geoapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	eventEndpoint             = "internal/v1/events/"
	unnotifiedEventsEndpoint  = "internal/v1/events/near"
	markEventNotifiedEndpoint = "/notified"
)

func (c *Client) PutEvent(ctx context.Context, obj EventInput) (Event, error) {
	const op = "geoapi.Client.PutEvent"

	var raw Event

	if err := c.put(ctx, eventEndpoint, obj, &raw); err != nil {
		return Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}

func (c *Client) UnnotifiedEventsNear(
	ctx context.Context, lat, lon float64, radius int, day time.Time,
) ([]Event, error) {
	const op = "geoapi.Client.UnnotifiedEventsNear"

	query := url.Values{
		"lat":    {strconv.FormatFloat(lat, 'f', -1, 64)},
		"lon":    {strconv.FormatFloat(lon, 'f', -1, 64)},
		"radius": {strconv.Itoa(radius)},
		"day":    {day.Format(time.DateOnly)},
	}

	var raw []Event

	if err := c.get(
		ctx, unnotifiedEventsEndpoint, query, &raw,
	); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}

func (c *Client) MarkEventNotified(ctx context.Context, ID string) error {
	const op = "geoapi.Client.MarkEventNotified"

	body := struct {
		ID string `json:"id"`
	}{ID: ID}

	endpoint := fmt.Sprintf("/internal/v1/events/%s/notified", ID)
	if err := c.put(ctx, endpoint, body, nil); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
