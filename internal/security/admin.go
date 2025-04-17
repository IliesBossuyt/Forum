package security

import (
	"Forum/internal/database"
	"encoding/json"
	"log"
	"net/http"
)

// Active ou désactive le bannissement d'un utilisateur
func ToggleBanUser(w http.ResponseWriter, r *http.Request) {
	// Vérifie la session de l'administrateur
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	// Vérifie les droits d'administration
	_, role, valid := ValidateSecureToken(cookie.Value, r.UserAgent())
	if !valid || (role != "admin" && role != "moderator") {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	// Décode les données de la requête
	var requestData struct {
		UserID string `json:"user_id"`
		Banned bool   `json:"banned"` // État à appliquer
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Si bannissement, supprime toutes les sessions
	if requestData.Banned {
		err := DeleteAllSessionsForUser(requestData.UserID)
		if err != nil {
			log.Println("Erreur suppression sessions:", err)
		}
	}

	// Met à jour l'état de bannissement
	_, err = database.DB.Exec("UPDATE users SET banned = ? WHERE id = ?", requestData.Banned, requestData.UserID)
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}

	// Retourne la confirmation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "État de bannissement mis à jour",
	})
}

// Supprime toutes les sessions d'un utilisateur
func DeleteAllSessionsForUser(userID string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}
