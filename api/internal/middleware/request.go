package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *loggingResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}

	w.status = status
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(body)
}

func (w *loggingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func LoggerMiddleware(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()

		writer := &loggingResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(writer, r)

		requestLogger := log.WithRequest(
			writer.status,
			r.Method,
			r.URL.Path,
			clientIP(r),
			r.UserAgent(),
			time.Since(startedAt),
		)

		switch {
		case writer.status >= http.StatusInternalServerError:
			requestLogger.Error("HTTP request completed")
		case writer.status >= http.StatusBadRequest:
			requestLogger.Warn("HTTP request completed")
		default:
			requestLogger.Info("HTTP request completed")
		}
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
