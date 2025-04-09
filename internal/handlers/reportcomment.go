package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

func ReportComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		CommentID int    `json:"comment_id"`
		Reason    string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(security.ContextUserIDKey).(string)

	err := models.CreateCommentReport(data.CommentID, userID, data.Reason)
	if err != nil {
		http.Error(w, "Erreur lors du signalement", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"success": true})
}
