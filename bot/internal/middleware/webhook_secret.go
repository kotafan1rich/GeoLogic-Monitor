package middleware

import "net/http"

const (
	BotApiSecretHeader = "X-Max-Bot-Api-Secret"
)

func WebhookSecret(expectedSecret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get(BotApiSecretHeader)
		if secret != expectedSecret {
			return
		}
		next.ServeHTTP(w, r)
	})
}
