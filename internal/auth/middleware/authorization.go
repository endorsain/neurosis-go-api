package middleware

import (
	"net/http"

	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
)

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := RolesFromContext(r.Context())
			if !ok {
				httptransport.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			for _, currentRole := range roles {
				if currentRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			httptransport.WriteError(w, http.StatusForbidden, "forbidden")
		})
	}
}
