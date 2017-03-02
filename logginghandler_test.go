package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		request            *http.Request
		userAgent          string
		forwardedFor       string
		expectedLogMessage string
	}{
		{
			request:            httptest.NewRequest("GET", "http://example.com/", nil),
			userAgent:          "User-Agent-Value",
			forwardedFor:       "X-Forwarded-For-Value",
			expectedLogMessage: `{"Date":"2000-01-01T00:00:00Z","RemoteAddress":"127.0.0.1","ForwardedFor":"X-Forwarded-For-Value","UserAgent":"User-Agent-Value","Method":"GET","URL":"http://example.com/"}` + "\n",
		},
		{
			request:            httptest.NewRequest("GET", "http://example.com/?test=1", nil),
			userAgent:          "",
			forwardedFor:       "",
			expectedLogMessage: `{"Date":"2000-01-01T00:00:00Z","RemoteAddress":"127.0.0.1","ForwardedFor":"","UserAgent":"","Method":"GET","URL":"http://example.com/?test=1"}` + "\n",
		},
	}

	for _, test := range tests {
		nextCalled := false
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)

		test.request.Header.Add("X-Forwarded-For", test.forwardedFor)
		test.request.Header.Add("User-Agent", test.userAgent)
		test.request.RemoteAddr = "127.0.0.1"

		h := LoggingHandler{
			Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
			}),
			Stdout: stdout,
			Stderr: stderr,
			Now:    func() time.Time { return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) },
		}

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, test.request)

		if !nextCalled {
			t.Errorf("Expected the next handler in the chain to be called.")
		}

		if stdout.String() != test.expectedLogMessage {
			t.Errorf("Expected log message '%v', but got '%v'", test.expectedLogMessage, stdout.String())
		}

		if stderr.String() != "" {
			t.Errorf("An unexpected message was written to stderr: %v", stderr.String())
		}
	}
}
