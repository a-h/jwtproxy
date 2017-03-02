package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	tests := []struct {
		request          *http.Request
		expectedResponse string
	}{
		{
			request:          httptest.NewRequest("GET", "/healthcheck", nil), // Access the healthcheck, for the OK message.
			expectedResponse: "OK",
		},
		{
			request:          httptest.NewRequest("GET", "/", nil), // Don't access the healthcheck.
			expectedResponse: "NextResponse",
		},
	}

	for _, test := range tests {
		h := HealthCheckHandler{
			Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("NextResponse"))
			}),
			Path: "/healthcheck",
		}

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, test.request)

		actualResponse := rec.Body.String()
		if actualResponse != test.expectedResponse {
			t.Errorf("Expected response '%v', but got '%v'", test.expectedResponse, actualResponse)
		}

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("Expected a 200 status code but got %v", rec.Result().StatusCode)
		}
	}
}
