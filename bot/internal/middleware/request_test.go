package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/logger"
)

func TestLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		wantLevel  string
		wantStatus string
	}{
		{name: "successful request", status: http.StatusOK, wantLevel: `"level":"INFO"`, wantStatus: `"status":200`},
		{name: "client error", status: http.StatusBadRequest, wantLevel: `"level":"WARN"`, wantStatus: `"status":400`},
		{name: "server error", status: http.StatusInternalServerError, wantLevel: `"level":"ERROR"`, wantStatus: `"status":500`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			log := logger.New(logger.LevelInfo, logger.FormatJson, false, &output)
			output.Reset()

			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
			})
			request := httptest.NewRequest(http.MethodPost, "/webhook", nil)
			request.RemoteAddr = "192.0.2.1:1234"
			request.Header.Set("User-Agent", "middleware-test")

			LoggerMiddleware(log, next).ServeHTTP(httptest.NewRecorder(), request)

			got := output.String()
			for _, want := range []string{
				tt.wantLevel,
				tt.wantStatus,
				`"msg":"HTTP request completed"`,
				`"method":"POST"`,
				`"path":"/webhook"`,
				`"ip":"192.0.2.1"`,
				`"user-agent":"middleware-test"`,
			} {
				if !strings.Contains(got, want) {
					t.Fatalf("log output %q does not contain %q", got, want)
				}
			}
		})
	}
}

func TestLoggerMiddlewareDefaultsStatusToOK(t *testing.T) {
	var output bytes.Buffer
	log := logger.New(logger.LevelInfo, logger.FormatJson, false, &output)
	output.Reset()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	LoggerMiddleware(log, next).ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/health", nil),
	)

	if got := output.String(); !strings.Contains(got, `"status":200`) {
		t.Fatalf("log output %q does not contain status 200", got)
	}
}
