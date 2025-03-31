package security

import (
	"Forum/internal/database"
	"encoding/json"
	"log"
	"net/http"
)

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
