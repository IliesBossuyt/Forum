package handlers

import (
	"encoding/json"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Handler pour liker/disliker un commentaire
func LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer userID depuis le contexte
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)

	var input struct {
		CommentID int `json:"comment_id"`
		Value     int `json:"value"` // 1 = like, -1 = dislike
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Appliquer le like/dislike
	err = models.ToggleCommentLike(userID, input.CommentID, input.Value)
	if err != nil {
		http.Error(w, "Erreur lors du like/dislike du commentaire", http.StatusInternalServerError)
		return
	}

	// Récupérer les nouveaux totaux
	likes, dislikes, err := models.GetCommentLikes(input.CommentID)
	if err != nil {
		http.Error(w, "Erreur récupération likes commentaires", http.StatusInternalServerError)
		return
	}

	// Réponse JSON
	response := map[string]interface{}{
		"success":  true,
		"likes":    likes,
		"dislikes": dislikes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
