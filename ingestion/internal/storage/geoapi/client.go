package geoapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v7"
)

type Client struct {
	httpClient     *http.Client
	baseURL        *url.URL
	attemptTimeout time.Duration
	maxRetries     uint
	newBackOff     func() backoff.BackOff
}

func New(hc *http.Client, baseURL string, attemptTimeout time.Duration, maxRetries uint) (*Client, error) {
	if hc == nil {
		return nil, ErrInvalidHTTPClient
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseURL, err)
	}

	if !validateURL(parsedURL) {
		return nil, fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	return &Client{
		httpClient:     hc,
		baseURL:        parsedURL,
		attemptTimeout: attemptTimeout,
		maxRetries:     maxRetries,
		newBackOff:     func() backoff.BackOff { return backoff.NewExponentialBackOff() },
	}, nil
}

func MustNew(hc *http.Client, baseURL string, attemptTimeout time.Duration, maxRetries uint) *Client {
	const op = "geoapi.MustNew"

	c, err := New(hc, baseURL, attemptTimeout, maxRetries)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init geo-api client: %v", op, err))
	}
	return c
}

func (c *Client) do(ctx context.Context, baseURL *url.URL, endpoint string, query url.Values) error {
	const op = "digitalspb.Client.do"

	target := baseURL.JoinPath(endpoint)
	if len(query) > 0 {
		target.RawQuery = query.Encode()
	}

	_, err := backoff.Retry(
		ctx,
		func() (struct{}, error) {
			return c.attempt(ctx, target.String())
		},
		backoff.WithBackOff(c.newBackOff()),
		backoff.WithMaxElapsedTime(0),
		backoff.WithMaxTries(c.maxRetries),
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}

	return nil
}

func (c *Client) attempt(ctx context.Context, target string) (struct{}, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, c.attemptTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(attemptCtx, http.MethodGet, target, nil)
	if err != nil {
		return struct{}{}, backoff.Permanent(fmt.Errorf("%w: %v", ErrBuildRequest, err))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return struct{}{}, backoff.Permanent(fmt.Errorf("%w: %v", ErrCreateRequest, err))
		}

		return struct{}{}, fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusOK:
	case resp.StatusCode == http.StatusBadRequest,
		resp.StatusCode == http.StatusUnauthorized,
		resp.StatusCode == http.StatusNotFound,
		resp.StatusCode == http.StatusInternalServerError:
		return struct{}{}, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	default:
		return struct{}{}, backoff.Permanent(fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode))
	}

	return struct{}{}, nil
}

func validateURL(u *url.URL) bool {
	return u != nil && u.Host != "" && u.Scheme != ""
}
