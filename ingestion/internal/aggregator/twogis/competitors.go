package twogis

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	placesEndpoint = "3.0/items"
	recordsPerPage = 10
	maxPages       = 5
	successCode    = 200

	spbCityID = "5348647327760881"
)

func (c *Client) ParseCompetitors(ctx context.Context, category string, since time.Time) ([]Item, error) {
	const op = "twogis.Client.ParseCompetitors"

	query := c.buildQuery(category, since)

	var (
		items      []Item
		totalPages = -1
	)

	for page := 1; ; page++ {
		query.Set("page", strconv.Itoa(page))

		var raw Response
		if err := c.getPlaces(ctx, query, &raw); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if len(raw.Result.Items) == 0 {
			break
		}

		if totalPages == -1 {
			totalPages = (raw.Result.Total + recordsPerPage - 1) / recordsPerPage
			items = make([]Item, 0, raw.Result.Total)
		}

		items = append(items, raw.Result.Items...)

		if page >= totalPages || page >= maxPages {
			break
		}
	}

	return items, nil
}

func (c *Client) getPlaces(ctx context.Context, query url.Values, out *Response) error {
	err := c.do(ctx, placesEndpoint, query, out)
	if err != nil {
		return err
	}

	if out.Meta.Code != successCode {
		if out.Meta.Error != nil && out.Meta.Error.Message != "" {
			return fmt.Errorf(
				"%w [%d]: %s", ErrUnexpectedProviderStatus, out.Meta.Code, out.Meta.Error.Message,
			)
		}
		return fmt.Errorf("%w [%d]", ErrUnexpectedProviderStatus, out.Meta.Code)
	}

	return nil
}
