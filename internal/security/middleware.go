package security

import (
	"context"
	"net/http"
)

type contextKey string

const (
	ContextUserIDKey contextKey = "userID"
	ContextRoleKey   contextKey = "role"
)

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID, role string
			role = "guest" // Valeur par défaut pour les visiteurs

			cookie, err := r.Cookie("session")
			if err == nil {
				// Token présent → validation
				uid, userRole, valid := ValidateSecureToken(cookie.Value, r.UserAgent())
				if valid {
					userID = uid
					role = userRole
				}
			}

			// ⛔ Si rôle non autorisé → redirection
			if !contains(allowedRoles, role) {
				http.Redirect(w, r, "/auth/unauthorized", http.StatusSeeOther)
				return
			}

			// ✅ Ajout des infos au contexte
			ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
			ctx = context.WithValue(ctx, ContextRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func contains(roles []string, target string) bool {
	for _, role := range roles {
		if role == target {
			return true
		}
	}
	return false
}
