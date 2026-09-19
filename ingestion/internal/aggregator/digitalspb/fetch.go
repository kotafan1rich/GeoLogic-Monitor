package digitalspb

import (
	"context"
	"fmt"
	"net/url"
)

const (
	egsGateV1 = 1
	egsGateV2 = 2
)

func fetchSpbClassifGate[T any](
	ctx context.Context, c *Client, baseURL *url.URL, endpoint, op string,
) ([]T, error) {
	if !validateURL(baseURL) {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidURL)
	}

	var raw SpbClassifGateResponse[T]

	err := c.do(ctx, baseURL, endpoint, url.Values{}, &raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return raw.Results, nil
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
