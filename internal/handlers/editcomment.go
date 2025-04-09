package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

func EditComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(security.ContextUserIDKey).(string)

	var input struct {
		CommentID int    `json:"comment_id"`
		Content   string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Content == "" {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Vérification d’auteur
	authorID, err := models.GetCommentAuthorID(input.CommentID)
	if err != nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
		return
	}

	if authorID != userID {
		http.Error(w, "Non autorisé", http.StatusForbidden)
		return
	}

	// Mise à jour
	err = models.UpdateCommentContent(input.CommentID, input.Content)
	if err != nil {
		http.Error(w, "Erreur lors de la modification", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
