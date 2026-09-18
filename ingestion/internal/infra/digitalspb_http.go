package infra

import (
	"net/http"
	"time"
)

func NewDigitalSpbHTTPClient(
	maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost int, reqTimeout time.Duration,
) *http.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.MaxIdleConns = maxIdleConns
	tr.MaxIdleConnsPerHost = maxIdleConnsPerHost
	tr.MaxConnsPerHost = maxConnsPerHost

	return &http.Client{
		Transport: tr,
		Timeout:   reqTimeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
