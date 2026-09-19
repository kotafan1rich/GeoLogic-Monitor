package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
)

func BotToken(expectedToken string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || !tokensEqual(parts[1], expectedToken) {
			response.WriteError(w, app.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IngestionToken(expectedToken string, next http.Handler) http.Handler {
	return BotToken(expectedToken, next)
}

func tokensEqual(actual, expected string) bool {
	if actual == "" || expected == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) == 1
}
