package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

// Gère la modification d'un commentaire
func EditComment(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupère l'ID de l'utilisateur
	userID := r.Context().Value(security.ContextUserIDKey).(string)

	// Décode le contenu et l'ID du commentaire depuis le JSON
	var input struct {
		CommentID int    `json:"comment_id"`
		Content   string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Content == "" {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Vérifie si l'utilisateur est l'auteur du commentaire
	authorID, err := models.GetCommentAuthorID(input.CommentID)
	if err != nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
		return
	}

	if authorID != userID {
		http.Error(w, "Non autorisé", http.StatusForbidden)
		return
	}

	// Met à jour le contenu du commentaire
	err = models.UpdateCommentContent(input.CommentID, input.Content)
	if err != nil {
		http.Error(w, "Erreur lors de la modification", http.StatusInternalServerError)
		return
	}

	// Retourne la confirmation de modification
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
