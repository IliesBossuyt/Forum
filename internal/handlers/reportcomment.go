package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

// Gère le signalement d'un commentaire
func ReportComment(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Structure pour décoder les données du signalement
	var data struct {
		CommentID int    `json:"comment_id"` // ID du commentaire signalé
		Reason    string `json:"reason"`     // Raison du signalement
	}

	// Décode les données de la requête
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Récupère l'ID de l'utilisateur qui signale
	userID := r.Context().Value(security.ContextUserIDKey).(string)

	// Crée le signalement dans la base de données
	err := models.CreateCommentReport(data.CommentID, userID, data.Reason)
	if err != nil {
		http.Error(w, "Erreur lors du signalement", http.StatusInternalServerError)
		return
	}

	// Répond avec succès
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}
