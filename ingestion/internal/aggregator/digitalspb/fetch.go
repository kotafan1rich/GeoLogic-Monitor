package digitalspb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/infra"
)

const (
	egsGateV1 = 1
	egsGateV2 = 2

	recordsPerPage = 500
)

func fetchSpbClassifGate[T any](
	ctx context.Context, c *Client, src a.URL, endpoint, op string,
) ([]T, error) {
	if err := src.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var results []T

	page := 1
	query := url.Values{
		"page":     {strconv.Itoa(page)},
		"per_page": {strconv.Itoa(recordsPerPage)},
	}

	for {
		var raw SpbClassifGateResponse[T]

		if err := c.do(ctx, src.URL(), endpoint, query, &raw); err != nil {
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
	ctx context.Context, c *Client, src a.URL, ver int, endpoint, op string,
) ([]T, error) {
	if err := src.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if ver == egsGateV1 {
		var raw EgsGateResponseV1[T]

		if err := c.do(ctx, src.URL(), endpoint, url.Values{}, &raw); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return raw.List, nil
	} else {
		var raw EgsGateResponseV2[T]

		if err := c.do(ctx, src.URL(), endpoint, url.Values{}, &raw); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return raw.Data, nil
	}
}

func fetchYazzhGate[T any](
	ctx context.Context, c *Client, src a.URL, endpoint, op string,
) ([]T, error) {
	if err := src.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var raw YazzhGateResponse[T]

	err := c.do(ctx, src.URL(), endpoint, url.Values{}, &raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return raw.Data, nil
}

func fetchStatic[T any](ctx context.Context, src a.File, op string) ([]T, error) {
	if err := src.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	f, err := os.Open(src.Path())
	if err != nil {
		return nil, fmt.Errorf("%s: %w [%s]: %v", op, ErrOpenFile, src, err)
	}
	defer f.Close()

	var raw []T

	if err := json.NewDecoder(io.LimitReader(f, infra.MaxBodySize)).Decode(&raw); err != nil {
		return nil, fmt.Errorf("%s: %w [%s]: %v", op, ErrDecodeData, src, err)
	}

	return raw, nil
}
