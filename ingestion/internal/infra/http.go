package infra

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPOptions struct {
	Name                string
	AuthToken           string
	RequireAuth         bool
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	RequestTimeout      time.Duration
	RateLimit           RateLimit
}

type authTransport struct {
	token string
	next  http.RoundTripper
}

func (a *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	reqClone := r.Clone(r.Context())
	reqClone.Header.Set("Authorization", "Bearer "+a.token)

	return a.next.RoundTrip(reqClone)
}

func NewHTTPClient(o HTTPOptions) (*http.Client, error) {
	if o.RequireAuth && o.AuthToken == "" {
		return nil, fmt.Errorf("%w [%s]", ErrInvalidAuthToken, o.Name)
	}

	if o.RequestTimeout <= 0 {
		return nil, fmt.Errorf("%w [%s]", ErrInvalidRequestTimeout, o.Name)
	}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.MaxIdleConns = o.MaxIdleConns
	tr.MaxIdleConnsPerHost = o.MaxIdleConnsPerHost
	tr.MaxConnsPerHost = o.MaxConnsPerHost

	var rt http.RoundTripper = tr

	if o.AuthToken != "" {
		rt = &authTransport{token: o.AuthToken, next: rt}
	}

	if o.RateLimit.Enabled() {
		rt = &rateLimitTransport{limiter: NewLimiter(o.RateLimit), next: rt}
	}

	return &http.Client{
		Transport: rt,
		Timeout:   o.RequestTimeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, nil
}

func MustNewHTTPClient(o HTTPOptions) *http.Client {
	const op = "infra.MustNewHTTPClient"

	hc, err := NewHTTPClient(o)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init http client [%s]: %v", op, o.Name, err))
	}

	return hc
}
