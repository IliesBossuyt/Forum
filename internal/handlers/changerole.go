package handlers

import (
	"Forum/internal/database"
	"encoding/json"
	"net/http"
)

// Rôles autorisés
var validRoles = map[string]bool{
	"user":      true,
	"moderator": true,
	"admin":     true,
}

func ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	// Lire le JSON envoyé par le frontend
	var requestData struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Vérifier si le rôle est valide
	if !validRoles[requestData.Role] {
		http.Error(w, "Rôle invalide", http.StatusBadRequest)
		return
	}

	// Mettre à jour le rôle dans la base de données
	_, err := database.DB.Exec("UPDATE users SET role = ? WHERE id = ?", requestData.Role, requestData.UserID)
	if err != nil {
		http.Error(w, "Erreur mise à jour", http.StatusInternalServerError)
		return
	}

	// Réponse JSON pour confirmer la mise à jour
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
