package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/cenkalti/backoff/v7"
)

const (
	MaxBodySize = 64 << 20

	defaultAttemptTimeout = 30 * time.Second

	mediaTypeJSON = "application/json"
	mediaTypeForm = "application/x-www-form-urlencoded"
)

type RetryPolicy func(status int) bool

func RetryTransient(status int) bool {
	switch {
	case status == http.StatusTooManyRequests, status == http.StatusRequestTimeout:
		return true
	case status >= http.StatusInternalServerError:
		return true
	default:
		return false
	}
}

type Requester struct {
	hc             *http.Client
	attemptTimeout time.Duration
	maxRetries     uint
	retry          RetryPolicy
	newBackOff     func() backoff.BackOff
}

func NewRequester(
	hc *http.Client, attemptTimeout time.Duration, maxRetries uint, retry RetryPolicy,
) (*Requester, error) {
	if hc == nil {
		return nil, ErrInvalidHTTPClient
	}

	if attemptTimeout <= 0 {
		attemptTimeout = defaultAttemptTimeout
	}

	if retry == nil {
		retry = RetryTransient
	}

	return &Requester{
		hc:             hc,
		attemptTimeout: attemptTimeout,
		maxRetries:     maxRetries,
		retry:          retry,
		newBackOff:     func() backoff.BackOff { return backoff.NewExponentialBackOff() },
	}, nil
}

func MustNewRequester(
	name string, hc *http.Client, attemptTimeout time.Duration, maxRetries uint, retry RetryPolicy,
) *Requester {
	const op = "infra.MustNewRequester"

	r, err := NewRequester(hc, attemptTimeout, maxRetries, retry)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init %q requester: %v", op, name, err))
	}

	return r
}

func Target(baseURL *url.URL, endpoint string, query url.Values) string {
	target := baseURL.JoinPath(endpoint)

	if len(query) > 0 {
		target.RawQuery = query.Encode()
	}

	return target.String()
}

func (r *Requester) JSON(ctx context.Context, method, target string, body, out any) error {
	var payload []byte

	if body != nil {
		var err error

		payload, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrMarshalData, err)
		}
	}

	return r.do(ctx, method, target, payload, mediaTypeJSON, out)
}

func (r *Requester) Form(ctx context.Context, method, target string, form url.Values, out any) error {
	return r.do(ctx, method, target, []byte(form.Encode()), mediaTypeForm, out)
}

func (r *Requester) do(
	ctx context.Context, method, target string, payload []byte, contentType string, out any,
) error {
	_, err := backoff.Retry(
		ctx,
		func() (struct{}, error) {
			return struct{}{}, r.attempt(ctx, method, target, payload, contentType, out)
		},
		backoff.WithBackOff(r.newBackOff()),
		backoff.WithMaxElapsedTime(0),
		backoff.WithMaxTries(r.maxRetries),
	)

	return err
}

func (r *Requester) attempt(
	ctx context.Context, method, target string, payload []byte, contentType string, out any,
) error {
	attemptCtx, cancel := context.WithTimeout(ctx, r.attemptTimeout)
	defer cancel()

	resetOut(out)

	var bodyReader io.Reader
	if payload != nil {
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(attemptCtx, method, target, bodyReader)
	if err != nil {
		return backoff.Permanent(fmt.Errorf("%w: %v", ErrBuildRequest, err))
	}

	if payload != nil {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := r.hc.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return backoff.Permanent(fmt.Errorf("%w: %v", ErrDoRequest, err))
		}

		return fmt.Errorf("%w: %v", ErrDoRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		drain(resp.Body)

		if r.retry(resp.StatusCode) {
			return fmt.Errorf("%w [%d]", ErrUnexpectedStatus, resp.StatusCode)
		}

		return backoff.Permanent(fmt.Errorf("%w [%d]", ErrUnexpectedStatus, resp.StatusCode))
	}

	if out == nil {
		drain(resp.Body)
		return nil
	}

	lr := &io.LimitedReader{R: resp.Body, N: MaxBodySize + 1}

	if err := json.NewDecoder(lr).Decode(out); err != nil {
		if lr.N <= 0 {
			return backoff.Permanent(fmt.Errorf("%w: %w", ErrDecodeResponse, ErrBodyTooLarge))
		}
		return fmt.Errorf("%w: %v", ErrDecodeResponse, err)
	}

	return nil
}

func resetOut(out any) {
	if out == nil {
		return
	}

	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return
	}

	elem := v.Elem()
	if elem.CanSet() {
		elem.Set(reflect.Zero(elem.Type()))
	}
}

func drain(r io.Reader) {
	_, _ = io.Copy(io.Discard, io.LimitReader(r, MaxBodySize+1))
}
