package infra

import (
	"fmt"
	"net/http"
	"time"
)

type authTransport struct {
	token string
	next  http.RoundTripper
}

func (a *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	reqClone := r.Clone(r.Context())
	reqClone.Header.Set("Authorization", "Bearer "+a.token)
	return a.next.RoundTrip(reqClone)
}

func NewGeoAPIHTTPClient(
	authToken string, maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost int, reqTimeout time.Duration,
) (*http.Client, error) {
	if authToken == "" {
		return nil, ErrInvalidAuthToken
	}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.MaxIdleConns = maxIdleConns
	tr.MaxIdleConnsPerHost = maxIdleConnsPerHost
	tr.MaxConnsPerHost = maxConnsPerHost

	var rt http.RoundTripper = tr
	rt = &authTransport{token: authToken, next: rt}

	return &http.Client{
		Transport: rt,
		Timeout:   reqTimeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, nil
}

func MustNewGeoAPIHTTPClient(
	authToken string, maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost int, reqTimeout time.Duration,
) *http.Client {
	const op = "infra.MustNewGeoAPIHTTPClient"

	hc, err := NewGeoAPIHTTPClient(
		authToken, maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost, reqTimeout,
	)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init geo-api http client: %v", op, err))
	}

	return hc
}
