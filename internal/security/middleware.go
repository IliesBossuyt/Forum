package security

import (
	"Forum/internal/models"
	"context"
	"net/http"
	"time"
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
			role = "guest" // Par défaut

			cookie, err := r.Cookie("session")
			if err == nil {
				// Vérifier le token (signature et User-Agent)
				uid, userRole, valid := ValidateSecureToken(cookie.Value, r.UserAgent())
				if valid {
					userID = uid
					role = userRole

					// Vérifier la session en base de données
					sessionUUID := ExtractUUID(cookie.Value)
					storedUserID, storedRole, expiresAt, err := GetUserIDFromSession(sessionUUID)
					if err != nil || storedUserID != userID || storedRole != role || expiresAt.Before(time.Now()) {
						http.Error(w, "Session invalide ou expirée", http.StatusUnauthorized)
						return
					}

					// Vérifier si l'utilisateur est banni
					user, err := models.GetUserByID(userID)
					if err != nil || user == nil || user.Banned {
						http.Error(w, "Utilisateur invalide ou banni", http.StatusUnauthorized)
						return
					}
				}
			}

			// Vérification du rôle autorisé
			if !contains(allowedRoles, role) {
				http.Error(w, "Accès refusé", http.StatusForbidden)
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			// Injecter l'utilisateur dans le contexte
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
