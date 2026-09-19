package digitalspb

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	egsGateV1 = 1
	egsGateV2 = 2

	recordsPerPage = 500
)

func fetchSpbClassifGate[T any](
	ctx context.Context, c *Client, baseURL *url.URL, endpoint, op string,
) ([]T, error) {
	if !validateURL(baseURL) {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidURL)
	}

	var results []T

	page := 1
	query := url.Values{
		"page":     {strconv.Itoa(page)},
		"per_page": {strconv.Itoa(recordsPerPage)},
	}

	for {
		var raw SpbClassifGateResponse[T]

		if err := c.do(ctx, baseURL, endpoint, query, &raw); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if results == nil {
			results = make([]T, 0, max(raw.Count, len(raw.Results)))
		}

		results = append(results, raw.Results...)

		if raw.Next == nil || *raw.Next == "" {
			return results, nil
		}

		page += 1
		query.Set("page", strconv.Itoa(page))
	}
}

func fetchEgsGate[T any](
	ctx context.Context, c *Client, baseURL *url.URL, ver int, endpoint, op string,
) ([]T, error) {
	if !validateURL(baseURL) {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidURL)
	}

	if ver == egsGateV1 {
		var raw EgsGateResponseV1[T]

		if err := c.do(ctx, baseURL, endpoint, url.Values{}, &raw); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return raw.List, nil
	} else {
		var raw EgsGateResponseV2[T]

		if err := c.do(ctx, baseURL, endpoint, url.Values{}, &raw); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return raw.Data, nil
	}
}

func fetchYazzhGate[T any](
	ctx context.Context, c *Client, baseURL *url.URL, endpoint, op string,
) ([]T, error) {
	if !validateURL(baseURL) {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidURL)
	}

	var raw YazzhGateResponse[T]

	err := c.do(ctx, baseURL, endpoint, url.Values{}, &raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return raw.Data, nil
}
