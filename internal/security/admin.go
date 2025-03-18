package security

import (
	"context"
	"net/http"
)

// Middleware pour restreindre l'accès aux admins uniquement
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer le cookie de session
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Non autorisé", http.StatusUnauthorized)
			return
		}

		// Vérifier le token et récupérer l'ID utilisateur + rôle
		userID, role, valid := ValidateSecureToken(cookie.Value, r.UserAgent())
		if !valid || role != "admin" {
			http.Error(w, "Accès refusé", http.StatusForbidden)
			return
		}

		// Ajouter l'ID utilisateur au contexte (utile pour d'autres handlers)

		type contextKey string

		const userIDKey = contextKey("userID")
		const roleKey = contextKey("role")
		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)

		// Continuer l'exécution de la route protégée
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
