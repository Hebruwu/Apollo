package middleware

import (
	"net/http"
)

// StrictAuth is a middleware that enforces strict authentication.
// This means that the JWT version must match the version stored in the database.
// Must be used for routes that modify user data.
func StrictAuth(next http.Handler) http.Handler {
	handlerFunction := func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement authentication logic
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFunction)
}

// QuickAuth is a middleware that enforces authentication.
// This means that the JWT version must match the version stored in the database.
// Must be used for routes that require authentication, but do not modify user data.
// If user data is modified, use StrictAuth instead.
func QuickAuth(next http.Handler) http.Handler {
	handlerFunction := func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement authentication logic
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFunction)
}
