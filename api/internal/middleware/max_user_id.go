package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
)

type contextKey string

const maxUserIDKey contextKey = "max_user_id"

func MaxUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawID := r.Header.Get("X-Max-User-Id")
		if rawID == "" {
			response.WriteError(w, app.Wrap(errors.New("X-Max-User-Id header is required"), app.ErrUnauthorized))
			return
		}

		maxUserID, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil || maxUserID <= 0 {
			response.WriteError(w, app.Wrap(errors.New("invalid X-Max-User-Id header"), app.ErrUnauthorized))
			return
		}

		ctx := context.WithValue(r.Context(), maxUserIDKey, maxUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetMaxUserID(ctx context.Context) (int64, bool) {
	maxUserID, ok := ctx.Value(maxUserIDKey).(int64)
	return maxUserID, ok
}
