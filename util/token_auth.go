package util

import (
	"fmt"
	"net/http"
)

// TokenAuthMiddleware makes a bearer token be required for all requests
type TokenAuthMiddleware struct {
	token string
	next  http.Handler
}

// NewTokenAuthMiddleware does what it says
func NewTokenAuthMiddleware(token string, next http.Handler) *TokenAuthMiddleware {
	return &TokenAuthMiddleware{token, next}
}

// ServeHTTP implements http.Handler
func (mw *TokenAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", mw.token) {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}

	mw.next.ServeHTTP(w, r)
}
