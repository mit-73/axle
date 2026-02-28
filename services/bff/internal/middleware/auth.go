package middleware

import (
	"net/http"
)

// Auth is a stub middleware — currently passes all requests through.
// Will be replaced with Ory integration when auth is wired up.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: validate Bearer token via Ory Kratos/Hydra session API.
		// On success: inject user claims into context.
		// On failure: http.Error(w, "Unauthorized", http.StatusUnauthorized)
		next.ServeHTTP(w, r)
	})
}
