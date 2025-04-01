package auth

import (
	"context"
	"net/http"

	"github.com/pdridh/service-needs-app/backend/api"
)

// Takes a handler function and only calls it if
// the jwt token it extracts from the request's is valid.
// The next handler function is called with the userid in context
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jCookie, err := r.Cookie("jwt")
		if err != nil {
			api.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
			return
		}

		j := jCookie.Value

		t, err := ValidateJWT(j)
		if err != nil {
			api.WriteError(w, r, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		c, err := UserClaimsFromJWT(t)
		if err != nil {
			api.WriteError(w, r, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		// If the token is valid and the claims were extracted then create a new context with the userid
		// For future handlers to access
		ctx := context.WithValue(r.Context(), api.ContextUserKey, c.UserID)

		next(w, r.WithContext(ctx))
	}
}
