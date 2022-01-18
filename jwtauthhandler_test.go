package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

import "net/http/httptest"

func TestJWTHandling(t *testing.T) {
	tests := []struct {
		name               string
		request            *http.Request
		expectedStatusCode int
		expectedBody       string
		expectedNextCalled bool
		now                func() time.Time
		authHeader         string
	}{
		{
			name:               "missing JWT header",
			request:            httptest.NewRequest("GET", "/", nil),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Required authorization token not found",
		},
		{
			name: "Junk format for authorization header",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "nonsense")
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Authorization header format must be Bearer {token}",
		},
		{
			name: "Invalid JWT",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "Bearer sfdjkfjk.sdsdads.asdasd")
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "invalid character",
		},
		{
			name: "Missing exp field",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSJ9.j2MxQgF6XbleT04c7M_zvP_JlL8M1Af5foUO5lf4UH5PLzpkjdZtVgbGkNUljQWZJ1PGe4VCwrPf_jJ2yai4hQk7xlOFGjnYJoF8D1wc1phPfAH_RDnIcIRmHbU-d3f5HsKyQI5v1Ba0dS4dzwxWyM0H80gKnhy-2afzWndsNVa0oFA9mrxUow7HveS77bRdoA5WzJPEQCaKg_2uDn3KnmI68qntaTdHJC5VULYclalsr3apxUXKNtMgBscq6r57CSqlOrnmCy9S9pGDvSFk4KbF4ZMhFYBB4YDC25vSFgWDg3XGev0_Tkw0crs0ndgXJPdykeKlLzQA3i1zE4_ysw")
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "token expired",
		},
		{
			name: "Expired on 1st January 2017",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSIsImV4cCI6MTQ4MzIyODgwMH0.uTEmiM8qDMPpKvhaNRrir7wQUJ0mtwMdafjg4y5pADBtyLr93JYhbSxsugIA0TqLLQmLgUyTZGVy3ACzbvHdsyCrSgdVqaN7zFuN3v2XrYzfHFwU4fZIB5ZnuK-V70TMOWYdSOOU-vNZts35FLWvdEM12u4vZrHkV8sBkSYuV137y3APse463ejBbqlLEljbZ2RXlJHwB8OMZHzDpPE-dpl9AkltAZp-xJM3EbD1HHIAzZcBUK1HhRe8-GBzKayEnBWN9vOfsOEwmO5o4RE5Tf7v1NkQChCQiwm-OXaQRHIem0I3v8R4YS4b1ErFcfbDbWy5FFY0jJqog8-3_5MqyQ")
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "token expired",
		},
		{
			name: "Expired on 1st January 2017, but today is 1st January 2016",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODMyMjg4MDB9.zt87SYVKKbewLboxen033MGGCyW9fqoivUpeWe7zSi6QorRBXsknsP9a6AABy7UrcOu0jMzmzoCNkDODsw3VE-4wLnnvga_NlRsy70Cib_2XIUVpy6yTm9-GR8dqTU0Lckxdeu2EIfh6FrZiGQ8KPHr23tgSpTysPKL5kmUweDgu58Rfl2qK0ImzuZ3L7PGsqeTcRyWMTdVPUr1eDW6o17l8qnqj0vOlLWNjoLC1TSLJiZ3GzqVy1hj2JnT6V5bWdCzULGG-099ip8XRYr_yWje0UmCTZgRe9bKweuGvtDdrV5JMzCPS-CK-VArJWdcjymZtT8B5QitYZhhpKQTiCQ")
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "iss not valid",
			now:                func() time.Time { return time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC) },
		},
		{
			name: "Wrong signing algorithm",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZXhwIjoxNDgzMjI4ODAwfQ.-x4JLnDdmmyENZoY_Cr3E8_aShD_PpWih5vI7EfRqOc")
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "invalid signing algorithm, expected \"RS256\"",
			now:                func() time.Time { return time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC) },
		},
		{
			name: "Valid exp, valid issuer",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSIsImV4cCI6MTUxNDc2NDgwMH0.aGNUa78lwsARfdoClbTYeWoFmPJoLpyLOJBBlUQnt-VXVwcn9x0mKCzP6bBoS8eU27-iE2dZXvCMYgwcITocH6EAP5MKgeARUGYaT30MvtokdLTYCXlgW1TQiT3QLNdae0wUSzTXgN6BkYckYeZqlyI77m15tJTMQCYkQOfEIPIUl80nwYOR1cPNheZ0tClYUUfqGG-QOcO9gEAN5C83lMdfikFoNfIXlVCwcDgf7iLll9VpGaKCEjZfKGoRkGO9VhsLgJgMZzLWJaPack25lkepc_jGKRcc4i8q_c9Um1Hzv4E8WKOg9DwgOgG7GY_rk7yXytya0ie5Wm-CO-oupg")
				return r
			}(),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "OK",
			expectedNextCalled: true,
			now:                func() time.Time { return time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC) },
		},
		{
			name: "Valid with other authorization header",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("X-Other-Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSIsImV4cCI6MTUxNDc2NDgwMH0.aGNUa78lwsARfdoClbTYeWoFmPJoLpyLOJBBlUQnt-VXVwcn9x0mKCzP6bBoS8eU27-iE2dZXvCMYgwcITocH6EAP5MKgeARUGYaT30MvtokdLTYCXlgW1TQiT3QLNdae0wUSzTXgN6BkYckYeZqlyI77m15tJTMQCYkQOfEIPIUl80nwYOR1cPNheZ0tClYUUfqGG-QOcO9gEAN5C83lMdfikFoNfIXlVCwcDgf7iLll9VpGaKCEjZfKGoRkGO9VhsLgJgMZzLWJaPack25lkepc_jGKRcc4i8q_c9Um1Hzv4E8WKOg9DwgOgG7GY_rk7yXytya0ie5Wm-CO-oupg")
				return r
			}(),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "OK",
			expectedNextCalled: true,
			now:                func() time.Time { return time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC) },
			authHeader:         "X-Other-Authorization",
		},
		{
			name: "Valid exp, issuer happens to be a number",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOjEyMywiZXhwIjoxNTE0NzY0ODAwfQ.zWTab-8hHyVomiSgEmEeyOvJkg3P4Gzvg_67s1ezMFHJTuLntrunM99eWCA4jymuzjoncrjlwWIsYRxWaQLgk49JHWIIqmmlZwYR65CdGPSWCBSZN7Piem0KmYgCd7jSK_W1la31gmewT2-CDDAXi0pjprcZ1X47M0lhGh43jHUa7IolMCAqH0qO9haR7HEBfTABlZcczESgeildtbEXR9hYdfqG9nvdSucGmM8TdiZRAd7qgKPSKeGLwN4KuU28jkdO8U-RavFpvx2ss1P5DwvCnz5G_DFHjGMQdVJFxh7jO6Zh_TfMcbpEZN1W-AyEFCYwqAHSQll9AAcRmikIEQ")
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "invalid claims",
			now:                func() time.Time { return time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC) },
		},
		{
			name: "Valid exp, unknown issuer",
			request: func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlMi5jb20iLCJleHAiOjE1MTQ3NjQ4MDB9.HQnOWB05T-kUz7O6GxbQ_I98XTyPBqsCjl0y3DEv1Sp_Z2wJS58LzXl5eckqxl0fK5i4J6v4v6ZDw07biy5T3OqnrLtYN4haycbg0OPRLIuC3UjcogqtsATvDhArDt_VWlUmb9RLpMEEGeB2uBepyMC3g_Wk_O6vkfWOeeM8zGSjOvL_JlaSLRhZo1RZMCXNUrPzHR3eON4fLVNEWfhS8W7WtBnOMWCC2jwvfROK9m7wblaDYwLzUEoghC_qvAC7D8Zly-zQ9Yos6TfgXxeoXM6bBAj9jHsoWR4dWO3WtvZEQeKXze9vvK8aOU5A8T9bFMEc0ul44D1B5okrg-7Veg")
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "iss not valid",
			now:                func() time.Time { return time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC) },
		},
	}

	keys := map[string]string{
		"example.com": `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA05IDL+Y6VaJvUWmI4vOH
G0mL3h8TqfQ/icg6PBiA01MPj/dHzM8mTbxsRlxEbtIHb82mOJeWavd+TmiLSPNX
pbcNu4ZoY+LCmpxf3C2Uk3kbL7APIOEw56QTDCH9znscRC4r75uXEfv38FCXySU+
uWmILAXXqEHHiFW2q4ieR6mHvR7qZ4gg3uJARsCMGkHMofOTtkVwjbh56lQboMWY
8vV0ap6fg7OuRjWt4RF5fd4kU3mWYLlJPnMqcjPifiCLzlqF4EP0lfcLRwHjMuD/
oFQers8auQMYKouhgqNuClBI4JZLznK9qULr5fuGjvJI5fS7UIY1yyvwx6NSlmSM
nQIDAQAB
-----END PUBLIC KEY-----`,
	}

	for _, test := range tests {
		actualNextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
			actualNextCalled = true
		})

		// Work out when to set the time
		now := time.Now
		if test.now != nil {
			now = test.now
		}

		authHeader := "Authorization"
		if test.authHeader != "" {
			authHeader = test.authHeader
		}

		handler := NewJWTAuthHandler(keys, now, authHeader, next)
		recorder := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(recorder, test.request)

		// Assert
		actual := recorder.Result()

		if test.expectedNextCalled != actualNextCalled {
			t.Errorf("expected next called %v, but got %v", test.expectedNextCalled, actualNextCalled)
		}

		if actual.StatusCode != test.expectedStatusCode {
			t.Errorf("expected status code %v, but got %v", test.expectedStatusCode, actual.StatusCode)
		}

		actualBody, err := ioutil.ReadAll(actual.Body)
		if err != nil {
			t.Errorf("failed to read body with error: %v", err)

		}
		if !strings.HasPrefix(string(actualBody), test.expectedBody) {
			t.Errorf("expected body to start with '%v' but got '%v'", test.expectedBody, string(actualBody))
		}
	}
}
