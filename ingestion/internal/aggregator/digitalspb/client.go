package digitalspb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v7"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
)

const maxBodySize = 64 << 20

type Client struct {
	httpClient     *http.Client
	BaseURLs       map[string]a.URL
	attemptTimeout time.Duration
	maxRetries     uint
	newBackOff     func() backoff.BackOff
}

func New(
	hc *http.Client, baseURLs map[string]string, attemptTimeout time.Duration, maxRetries uint,
) (*Client, error) {
	if hc == nil {
		return nil, ErrInvalidHTTPClient
	}

	if len(baseURLs) == 0 {
		return nil, ErrInvalidURLMap
	}

	parsedBaseURLs := make(map[string]a.URL, len(baseURLs))

	for source, raw := range baseURLs {
		src, err := a.NewURL(raw)
		if err != nil {
			return nil, fmt.Errorf("%w [%s]: %v", ErrInvalidURL, source, err)
		}

		parsedBaseURLs[source] = src
	}

	return &Client{
		httpClient:     hc,
		BaseURLs:       parsedBaseURLs,
		attemptTimeout: attemptTimeout,
		maxRetries:     maxRetries,
		newBackOff:     func() backoff.BackOff { return backoff.NewExponentialBackOff() },
	}, nil
}

func MustNew(
	hc *http.Client, baseURLs map[string]string, attemptTimeout time.Duration, maxRetries uint,
) *Client {
	const op = "digitalspb.MustNew"

	c, err := New(hc, baseURLs, attemptTimeout, maxRetries)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init digital-spb client: %v", op, err))
	}
	return c
}

func (c *Client) do(ctx context.Context, baseURL *url.URL, endpoint string, query url.Values, out any) error {
	const op = "digitalspb.Client.do"

	target := baseURL.JoinPath(endpoint)
	if len(query) > 0 {
		target.RawQuery = query.Encode()
	}

	bytes, err := backoff.Retry(
		ctx,
		func() ([]byte, error) {
			return c.attempt(ctx, target.String())
		},
		backoff.WithBackOff(c.newBackOff()),
		backoff.WithMaxElapsedTime(0),
		backoff.WithMaxTries(c.maxRetries),
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrUnmarshalData, err)
	}

	return nil
}

func (c *Client) attempt(ctx context.Context, target string) ([]byte, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, c.attemptTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(attemptCtx, http.MethodGet, target, nil)
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("%w: %v", ErrBuildRequest, err))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, backoff.Permanent(fmt.Errorf("%w: %v", ErrCreateRequest, err))
		}

		return nil, fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize+1))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadResponse, err)
	}

	switch {
	case resp.StatusCode == http.StatusOK:
	case resp.StatusCode == http.StatusTooManyRequests,
		resp.StatusCode == http.StatusRequestTimeout,
		resp.StatusCode >= 500 && resp.StatusCode <= 599:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	default:
		return nil, backoff.Permanent(fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode))
	}

	if len(bytes) > maxBodySize {
		return nil, fmt.Errorf("%w: body exceeds limit", ErrReadResponse)
	}

	if len(bytes) == 0 {
		return nil, fmt.Errorf("%w: empty response body", ErrReadResponse)
	}

	return bytes, nil
}
