package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/response"
)

const maxInitDataHeader = "X-Max-Init-Data"

var errInvalidInitData = errors.New("invalid MAX init data")

type maxInitDataUser struct {
	ID int64 `json:"id"`
}

func MiniAppInitData(botToken string, maxAge time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maxUserID, err := validateMAXInitData(
			r.Header.Get(maxInitDataHeader),
			botToken,
			maxAge,
			time.Now(),
		)
		if err != nil {
			response.WriteError(w, app.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), maxUserIDKey, maxUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateMAXInitData(
	rawInitData string,
	botToken string,
	maxAge time.Duration,
	now time.Time,
) (int64, error) {
	if rawInitData == "" || botToken == "" || maxAge <= 0 {
		return 0, errInvalidInitData
	}

	params, err := url.ParseQuery(rawInitData)
	if err != nil || len(params) == 0 {
		return 0, errInvalidInitData
	}
	for key, values := range params {
		if key == "" || len(values) != 1 {
			return 0, errInvalidInitData
		}
	}

	receivedHash, err := decodeInitDataHash(params.Get("hash"))
	if err != nil {
		return 0, errInvalidInitData
	}

	launchParams := buildLaunchParams(params)
	secretKey := calculateHMAC([]byte("WebAppData"), []byte(botToken))
	expectedHash := calculateHMAC(secretKey, []byte(launchParams))
	if !hmac.Equal(receivedHash, expectedHash) {
		return 0, errInvalidInitData
	}

	authDateUnix, err := strconv.ParseInt(params.Get("auth_date"), 10, 64)
	if err != nil || authDateUnix <= 0 {
		return 0, errInvalidInitData
	}
	authDate := time.Unix(authDateUnix, 0)
	if authDate.After(now) || now.Sub(authDate) > maxAge {
		return 0, errInvalidInitData
	}

	var user maxInitDataUser
	if err := json.Unmarshal([]byte(params.Get("user")), &user); err != nil || user.ID <= 0 {
		return 0, errInvalidInitData
	}

	return user.ID, nil
}

func decodeInitDataHash(value string) ([]byte, error) {
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != sha256.Size {
		return nil, errInvalidInitData
	}
	return decoded, nil
}

func buildLaunchParams(params url.Values) string {
	keys := make([]string, 0, len(params)-1)
	for key := range params {
		if key != "hash" {
			keys = append(keys, key)
		}
	}
	slices.Sort(keys)

	var result strings.Builder
	for i, key := range keys {
		if i > 0 {
			result.WriteByte('\n')
		}
		result.WriteString(key)
		result.WriteByte('=')
		result.WriteString(params.Get(key))
	}
	return result.String()
}

func calculateHMAC(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(data)
	return mac.Sum(nil)
}
