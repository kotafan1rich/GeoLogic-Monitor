package request

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{name: "valid body", body: `{"value":"ok"}`},
		{name: "unknown field", body: `{"unknown":"value"}`, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			var destination struct {
				Value string `json:"value"`
			}

			err := DecodeJSON(request, &destination)
			if (err != nil) != test.wantErr {
				t.Fatalf("DecodeJSON() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
