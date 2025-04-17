package handlers

import (
	"Forum/internal/database"
	"Forum/internal/models"
	"encoding/json"
	"net/http"
)

// Liste des rôles
var validRoles = map[string]bool{
	"user":      true,
	"moderator": true,
	"admin":     true,
}

// Change le rôle d'un utilisateur
func ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur et le nouveau rôle
	var requestData struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Vérifie si le rôle est valide
	if !validRoles[requestData.Role] {
		http.Error(w, "Rôle invalide", http.StatusBadRequest)
		return
	}

	// Met à jour le rôle dans la base de données
	_, err := database.DB.Exec("UPDATE users SET role = ? WHERE id = ?", requestData.Role, requestData.UserID)
	if err != nil {
		http.Error(w, "Erreur mise à jour", http.StatusInternalServerError)
		return
	}

	// Met à jour le rôle dans la session
	err = models.UpdateUserSessionRole(requestData.UserID, requestData.Role)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour de la session", http.StatusInternalServerError)
		return
	}

	// Retourne la confirmation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
