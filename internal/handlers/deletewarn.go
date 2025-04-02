package handlers

import (
	"Forum/internal/models"
	"encoding/json"
	"net/http"
)

func DeleteWarn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		WarnID int `json:"warn_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	err := models.DeleteWarnByID(input.WarnID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"success": true})
}
