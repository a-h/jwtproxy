package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

// NewLoggingHandler creates a handler which prints to stdout.
func NewLoggingHandler(next http.Handler) LoggingHandler {
	return LoggingHandler{
		Next:   next,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Now:    time.Now,
	}
}

// LoggingHandler logs the incoming requests to stdout.
type LoggingHandler struct {
	Next   http.Handler
	Stdout io.Writer
	Stderr io.Writer
	Now    clock
}

type clock func() time.Time

func (lh LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logEntry := struct {
		Date          time.Time
		RemoteAddress string
		ForwardedFor  string
		UserAgent     string
		Method        string
		URL           string
	}{
		Date:          lh.Now(),
		RemoteAddress: r.RemoteAddr,
		ForwardedFor:  r.Header.Get("X-Forwarded-For"),
		UserAgent:     r.UserAgent(),
		Method:        r.Method,
		URL:           r.URL.String(),
	}
	bytes, err := json.Marshal(logEntry)
	if err != nil {
		lh.Stderr.Write([]byte(err.Error()))
		lh.Stderr.Write([]byte("\n"))
	}
	lh.Stdout.Write(bytes)
	lh.Stdout.Write([]byte("\n"))
	lh.Next.ServeHTTP(w, r)
}
