package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBotToken(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
		expectedToken string
		wantStatus    int
		wantCalled    bool
	}{
		{
			name:          "valid token",
			authorization: "Bearer secret",
			expectedToken: "secret",
			wantStatus:    http.StatusNoContent,
			wantCalled:    true,
		},
		{
			name:          "missing header",
			expectedToken: "secret",
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "invalid token",
			authorization: "Bearer wrong",
			expectedToken: "secret",
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "invalid scheme",
			authorization: "Basic secret",
			expectedToken: "secret",
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "empty configured token",
			authorization: "Bearer secret",
			wantStatus:    http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusNoContent)
			})
			handler := BotToken(tt.expectedToken, next)
			request := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", nil)
			if tt.authorization != "" {
				request.Header.Set("Authorization", tt.authorization)
			}
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", responseRecorder.Code, tt.wantStatus)
			}
			if called != tt.wantCalled {
				t.Fatalf("next called = %t, want %t", called, tt.wantCalled)
			}
		})
	}
}
