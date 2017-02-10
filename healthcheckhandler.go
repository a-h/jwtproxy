package main

import "net/http"

// HealthCheckHandler returns HTTP 200 and 'OK' when hit, passing non-matching requests
// through to the next handler.
type HealthCheckHandler struct {
	Path string
	Next http.Handler
}

func (h HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == h.Path {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	h.Next.ServeHTTP(w, r)
}
