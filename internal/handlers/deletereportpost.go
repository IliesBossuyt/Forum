package handlers

import (
	"Forum/internal/models"
	"encoding/json"
	"net/http"
)

// Structure de la requête JSON
type deleteReportRequest struct {
	ReportID int `json:"report_id"`
}

func DeleteReportPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req deleteReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	err := models.DeleteReportByID(req.ReportID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression du signalement", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
