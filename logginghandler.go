package main

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingHandler logs the incoming requests to stdout.
type LoggingHandler struct {
	Next http.Handler
}

func (lh LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Date
	// IP address
	// Forwarded IP address
	// User agent
	// Verb
	// Path and Query
	fmt.Printf("%v %v %v %v %v %v\n", time.Now(), r.RemoteAddr, r.Header.Get("X-Forwarded-For"), r.UserAgent(), r.Method, r.URL)

	lh.Next.ServeHTTP(w, r)
}
