package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	myhttp "github.com/kimkachailo/auth-service/internal/http"
)

func TestHealthRoute(t *testing.T) {
	mux := myhttp.NewRouter()

	tests := []struct {
		name   string
		method string
		path   string
		want   int
	}{
		{"health", http.MethodGet, "/health", http.StatusOK},
		{"wrong method", http.MethodPost, "/health", http.StatusMethodNotAllowed},
		{"wrong path", http.MethodGet, "/wrong", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.want {
				t.Fatalf("status: got %d; want %d", res.StatusCode, tt.want)
			}
			if tt.want == http.StatusOK {
				if got := res.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf("content type: got %q", got)
				}
				var body struct {
					Status string `json:"status"`
				}
				if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body.Status != "ok" {
					t.Errorf("status field: got %q; want ok", body.Status)
				}
			}

		})
	}
}
