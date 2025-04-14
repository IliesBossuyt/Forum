package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

func LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupération de l'ID utilisateur depuis le contexte
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)

	// Structure d'entrée attendue
	var input struct {
		CommentID int `json:"comment_id"`
		Value     int `json:"value"` // 1 pour like, -1 pour dislike
	}

	// Décodage du corps JSON
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Appliquer le like/dislike d'abord
	added, wasLike, err := models.ToggleCommentLike(userID, input.CommentID, input.Value)
	if err != nil {
		http.Error(w, "Erreur lors du traitement du like", http.StatusInternalServerError)
		return
	}

	// Récupérer les nouveaux likes/dislikes
	likes, dislikes, err := models.GetCommentLikes(input.CommentID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des likes", http.StatusInternalServerError)
		return
	}

	if added {
		comment, err := models.GetCommentByID(input.CommentID)
		if err == nil && comment.UserID != userID {
			notifType := "dislike_comment"
			if wasLike {
				notifType = "like_comment"
			}

			models.CreateNotification(models.Notification{
				RecipientID: comment.UserID,
				SenderID:    userID,
				Type:        notifType,
				CommentID:   &input.CommentID,
			})
		}
	}

	// Réponse JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"likes":    likes,
		"dislikes": dislikes,
	})
}
