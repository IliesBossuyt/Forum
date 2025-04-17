package handlers

import (
	"Forum/internal/models"
	"encoding/json"
	"net/http"
)

// Structure de la requête JSON pour la suppression d'un signalement
type deleteReportRequest struct {
	ReportID int `json:"report_id"`
}

// Gère la suppression d'un signalement de post
func DeleteReportPost(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Décode l'ID du signalement depuis le JSON
	var req deleteReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Supprime le signalement de la base de données
	err := models.DeleteReportByID(req.ReportID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression du signalement", http.StatusInternalServerError)
		return
	}

	// Retourne la confirmation de suppression
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
