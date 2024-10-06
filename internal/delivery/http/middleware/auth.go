package middleware

import (
	"context"
	"net/http"

	"github.com/Karzoug/meower-post-service/internal/delivery/http/gen"
)

type authKey struct{}

var AuthUsernameKey authKey

// Auth is a middleware that adds the username
// from the JWT token in the request Authorization Header to the request context.
// (!) The middleware doesn't check if the token is valid - it's up to the gateway.
func Auth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// if spec not claim authentification
			if r.Context().Value(gen.BearerAuthScopes) == nil {
				next.ServeHTTP(w, r)
				return
			}

			// otherwise
			sub := r.Header.Get("X-User")
			if sub == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthUsernameKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
