package main

import "testing"
import "net/http"
import "net/url"
import "net/http/httptest"

func TestRewriteHandler(t *testing.T) {
	tests := []struct {
		IncomingURL         string
		Prefix              string
		ExpectedOutgoingURL string
	}{
		{
			IncomingURL:         "/api/user?id=123",
			Prefix:              "/api",
			ExpectedOutgoingURL: "/user?id=123",
		},
		{
			IncomingURL:         "/api/user?id=123",
			Prefix:              "/api/",
			ExpectedOutgoingURL: "/user?id=123",
		},
		{
			IncomingURL:         "/api/user?id=123&test=t",
			Prefix:              "",
			ExpectedOutgoingURL: "/api/user?id=123&test=t",
		},
	}

	for i, test := range tests {
		var actualURL *url.URL
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualURL = r.URL
		})

		rh := NewRewriteHandler(test.Prefix, testHandler)

		r := httptest.NewRequest("GET", test.IncomingURL, nil)
		w := httptest.NewRecorder()
		rh.ServeHTTP(w, r)

		if actualURL.String() != test.ExpectedOutgoingURL {
			t.Errorf("%d: expected '%s', got '%s'", i, test.ExpectedOutgoingURL, actualURL.String())
		}
	}
}
