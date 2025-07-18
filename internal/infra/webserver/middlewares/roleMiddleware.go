package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
)

func RoleMiddleware(roleDB database.RoleInterface, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			userRole, ok := claims["role"].(string)
			if !ok {
				claimsJSON, _ := json.Marshal(claims) // Converte claims para JSON
				http.Error(w, fmt.Sprintf(`{"error": "invalid token", "claims": %s}`, claimsJSON), http.StatusForbidden)
				return
			}

			// Primeiro tenta comparar diretamente o ID com o nome permitido
			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Caso não seja um nome, tenta buscar o nome da role pelo ID
			roleEntity, err := roleDB.FindRoleByID(userRole)
			if err == nil {
				for _, role := range allowedRoles {
					if roleEntity.Name == role {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
		})
	}
}
