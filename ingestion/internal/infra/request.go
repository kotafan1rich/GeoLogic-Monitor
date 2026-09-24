package infra

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

const (
	// MaxBodySize ограничивает объём вычитываемого тела ответа.
	MaxBodySize = 64 << 20

	defaultAttemptTimeout = 30 * time.Second
)

// RetryPolicy сообщает, имеет ли смысл повторить запрос при таком статусе.
type RetryPolicy func(status int) bool

// RetryTransient повторяет запрос при троттлинге, таймауте и ошибках сервера.
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

// Requester выполняет JSON-запросы с потайм-аутом на попытку и экспоненциальным
// backoff. Один и тот же Requester используют все исходящие клиенты сервиса,
// поэтому логика повторов, лимитов тела и обработки статусов живёт в одном месте.
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

// Target собирает адрес запроса из базового URL, эндпоинта и query.
func Target(baseURL *url.URL, endpoint string, query url.Values) string {
	target := baseURL.JoinPath(endpoint)

	if len(query) > 0 {
		target.RawQuery = query.Encode()
	}

	return target.String()
}

// JSON выполняет запрос, повторяя временные сбои, и разбирает ответ в out.
// Пустой body означает запрос без тела, пустой out — что тело ответа
// только проверяется, но не разбирается.
func (r *Requester) JSON(ctx context.Context, method, target string, body, out any) error {
	var payload []byte

	if body != nil {
		var err error

		payload, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrMarshalData, err)
		}
	}

	raw, err := backoff.Retry(
		ctx,
		func() ([]byte, error) { return r.attempt(ctx, method, target, payload) },
		backoff.WithBackOff(r.newBackOff()),
		backoff.WithMaxElapsedTime(0),
		backoff.WithMaxTries(r.maxRetries),
	)
	if err != nil {
		return err
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalData, err)
	}

	return nil
}

func (r *Requester) attempt(ctx context.Context, method, target string, payload []byte) ([]byte, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, r.attemptTimeout)
	defer cancel()

	var bodyReader io.Reader
	if payload != nil {
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(attemptCtx, method, target, bodyReader)
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("%w: %v", ErrBuildRequest, err))
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := r.hc.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, backoff.Permanent(fmt.Errorf("%w: %v", ErrDoRequest, err))
		}

		return nil, fmt.Errorf("%w: %v", ErrDoRequest, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, MaxBodySize+1))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadResponse, err)
	}

	if resp.StatusCode != http.StatusOK {
		statusErr := fmt.Errorf("%w: %d [%s]", ErrUnexpectedStatus, resp.StatusCode, target)

		if r.retry(resp.StatusCode) {
			return nil, statusErr
		}

		return nil, backoff.Permanent(statusErr)
	}

	if len(raw) > MaxBodySize {
		return nil, backoff.Permanent(fmt.Errorf("%w: %v", ErrReadResponse, ErrBodyTooLarge))
	}

	if len(raw) == 0 {
		return nil, fmt.Errorf("%w: %v", ErrReadResponse, ErrEmptyResponse)
	}

	return raw, nil
}
