package geoapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v7"
)

const maxBodySize = 64 << 20

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

func (c *Client) do(ctx context.Context, baseURL *url.URL, endpoint string, body, out any) error {
	target := baseURL.JoinPath(endpoint)

	bytes, err := backoff.Retry(
		ctx,
		func() ([]byte, error) {
			return c.attempt(ctx, target.String(), body)
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
		return fmt.Errorf("%w: %v", ErrUnmarshalData, err)
	}

	return nil
}

func (c *Client) attempt(ctx context.Context, target string, body any) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("%w: %v", ErrMarshalData, err))
	}

	attemptCtx, cancel := context.WithTimeout(ctx, c.attemptTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(attemptCtx, http.MethodPut, target, bytes.NewReader(payload))
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("%w: %v", ErrBuildRequest, err))
	}

	req.Header.Set("Content-Type", "application/json")

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

	// if resp.StatusCode != http.StatusOK {
	// 	var apiErr struct {
	// 		Code    string `json:"code"`
	// 		Message string `json:"message"`
	// 	}
	// 	_ = json.Unmarshal(bytes, &apiErr)

	// 	return nil, backoff.Permanent(fmt.Errorf("%w: %d %s: %s.\nPayload: %s",
	// 		ErrUnexpectedStatus, resp.StatusCode, apiErr.Code, apiErr.Message, string(payload)))
	// }

	switch {
	case resp.StatusCode == http.StatusOK:
	case resp.StatusCode == http.StatusBadRequest,
		resp.StatusCode == http.StatusUnauthorized,
		resp.StatusCode == http.StatusNotFound,
		resp.StatusCode == http.StatusInternalServerError:

		return nil, backoff.Permanent(fmt.Errorf("%w: %d. Payload: %s", ErrUnexpectedStatus, resp.StatusCode, string(payload)))
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

func validateURL(u *url.URL) bool {
	return u != nil && u.Host != "" && u.Scheme != ""
}
