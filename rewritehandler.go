package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// RewriteHandler takes the input URL and strips the prefix from it.
type RewriteHandler struct {
	PrefixToRemove string
	Next           http.Handler
}

// NewRewriteHandler creates a handler which strips the prefix from an incoming URL.
func NewRewriteHandler(prefixToRemove string, next http.Handler) RewriteHandler {
	return RewriteHandler{
		PrefixToRemove: prefixToRemove,
		Next:           next,
	}
}

func (h RewriteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.String()
	trimmedURL := strings.TrimPrefix(originalURL, h.PrefixToRemove)
	if !strings.HasPrefix(trimmedURL, "/") {
		trimmedURL = "/" + trimmedURL
	}
	updatedURL, err := url.Parse(trimmedURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error trimming prefix '%v' from '%v', could not parse resulting URL of '%v'", h.PrefixToRemove, originalURL, trimmedURL), http.StatusInternalServerError)
		return
	}
	r.URL = updatedURL
	h.Next.ServeHTTP(w, r)
}
