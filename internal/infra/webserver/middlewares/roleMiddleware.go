package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			userRole, ok := claims["role"].(string)
			if !ok {
				claimsJSON, _ := json.Marshal(claims) // Converte claims para JSON
				http.Error(w, fmt.Sprintf(`{"error": "invalid token", "claims": %s}`, claimsJSON), http.StatusForbidden)
				return
			}

			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
		})
	}
}
