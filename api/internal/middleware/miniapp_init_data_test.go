package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

const testMaxBotToken = "max-bot-token"

func TestValidateMAXInitData(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 20, 17, 0, 0, 0, time.UTC)
	valid := signTestInitData(t, testMaxBotToken, now, `{"id":67890,"first_name":"Max"}`)

	tests := []struct {
		name        string
		initData    string
		wantUserID  int64
		wantInvalid bool
	}{
		{
			name:       "valid init data",
			initData:   valid,
			wantUserID: 67890,
		},
		{
			name:        "missing init data",
			wantInvalid: true,
		},
		{
			name:        "tampered user",
			initData:    replaceTestParam(t, valid, "user", `{"id":1}`),
			wantInvalid: true,
		},
		{
			name:        "duplicate parameter",
			initData:    valid + "&auth_date=" + strconv.FormatInt(now.Unix(), 10),
			wantInvalid: true,
		},
		{
			name:        "expired init data",
			initData:    signTestInitData(t, testMaxBotToken, now.Add(-2*time.Hour), `{"id":67890}`),
			wantInvalid: true,
		},
		{
			name:        "future init data",
			initData:    signTestInitData(t, testMaxBotToken, now.Add(time.Minute), `{"id":67890}`),
			wantInvalid: true,
		},
		{
			name:        "invalid user",
			initData:    signTestInitData(t, testMaxBotToken, now, `{"id":0}`),
			wantInvalid: true,
		},
		{
			name:        "invalid encoding",
			initData:    "user=%zz",
			wantInvalid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userID, err := validateMAXInitData(tt.initData, testMaxBotToken, time.Hour, now)
			if tt.wantInvalid {
				if err == nil {
					t.Fatalf("validateMAXInitData returned user ID %d, want an error", userID)
				}
				return
			}
			if err != nil {
				t.Fatalf("validateMAXInitData returned an error: %v", err)
			}
			if userID != tt.wantUserID {
				t.Fatalf("user ID: got %d, want %d", userID, tt.wantUserID)
			}
		})
	}
}

func TestMiniAppInitDataAddsUserIDToContext(t *testing.T) {
	t.Parallel()

	initData := signTestInitData(t, testMaxBotToken, time.Now(), `{"id":67890}`)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetMaxUserID(r.Context())
		if !ok || userID != 67890 {
			t.Fatalf("MAX user ID: got %d, %t; want 67890, true", userID, ok)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/api/v1/tracked-locations", nil)
	request.Header.Set(maxInitDataHeader, initData)
	responseRecorder := httptest.NewRecorder()

	MiniAppInitData(testMaxBotToken, time.Hour, next).ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want %d", responseRecorder.Code, http.StatusNoContent)
	}
}

func TestMiniAppInitDataRejectsInvalidData(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	})
	request := httptest.NewRequest(http.MethodGet, "/api/v1/tracked-locations", nil)
	request.Header.Set(maxInitDataHeader, "invalid")
	responseRecorder := httptest.NewRecorder()

	MiniAppInitData(testMaxBotToken, time.Hour, next).ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want %d", responseRecorder.Code, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("next handler was called")
	}
}

func signTestInitData(t *testing.T, botToken string, authDate time.Time, user string) string {
	t.Helper()

	params := url.Values{
		"auth_date": {strconv.FormatInt(authDate.Unix(), 10)},
		"chat":      {`{"id":12345,"type":"DIALOG"}`},
		"query_id":  {"4c0ab423-342b-4e45-aea4-2747dbc500cd"},
		"user":      {user},
	}
	launchParams := "auth_date=" + params.Get("auth_date") +
		"\nchat=" + params.Get("chat") +
		"\nquery_id=" + params.Get("query_id") +
		"\nuser=" + params.Get("user")

	secretMAC := hmac.New(sha256.New, []byte("WebAppData"))
	_, _ = secretMAC.Write([]byte(botToken))
	signatureMAC := hmac.New(sha256.New, secretMAC.Sum(nil))
	_, _ = signatureMAC.Write([]byte(launchParams))
	params.Set("hash", hex.EncodeToString(signatureMAC.Sum(nil)))

	return params.Encode()
}

func replaceTestParam(t *testing.T, initData, key, value string) string {
	t.Helper()

	params, err := url.ParseQuery(initData)
	if err != nil {
		t.Fatalf("parse test init data: %v", err)
	}
	params.Set(key, value)
	return params.Encode()
}
