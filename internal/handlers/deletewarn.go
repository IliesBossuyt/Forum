package handlers

import (
	"Forum/internal/models"
	"encoding/json"
	"net/http"
)

// Gère la suppression d'un avertissement
func DeleteWarn(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Décode l'ID de l'avertissement depuis le JSON
	var input struct {
		WarnID int `json:"warn_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Supprime l'avertissement de la base de données
	err := models.DeleteWarnByID(input.WarnID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	// Retourne la confirmation de suppression
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}
