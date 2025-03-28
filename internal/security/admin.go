package security

import (
	"Forum/internal/database"
	"context"
	"encoding/json"
	"log"
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
		if !valid || (role != "admin" && role != "moderator") {
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

func ToggleBanUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	_, role, valid := ValidateSecureToken(cookie.Value, r.UserAgent())
	if !valid || (role != "admin" && role != "moderator") {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	var requestData struct {
		UserID string `json:"user_id"`
		Banned bool   `json:"banned"` // état à appliquer
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	if requestData.Banned {
		err := DeleteAllSessionsForUser(requestData.UserID)
		if err != nil {
			log.Println("Erreur suppression sessions:", err)
		}
	}

	_, err = database.DB.Exec("UPDATE users SET banned = ? WHERE id = ?", requestData.Banned, requestData.UserID)
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "État de bannissement mis à jour",
	})
}

func DeleteAllSessionsForUser(userID string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}
