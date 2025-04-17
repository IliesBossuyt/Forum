package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

// Gère les likes et dislikes sur les commentaires
func LikeComment(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupère l'ID de l'utilisateur
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)

	// Structure pour le like/dislike
	var input struct {
		CommentID int `json:"comment_id"`
		Value     int `json:"value"` // 1 pour like, -1 pour dislike
	}

	// Décode les données de la requête
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Applique le like/dislike et récupère l'état
	added, wasLike, err := models.ToggleCommentLike(userID, input.CommentID, input.Value)
	if err != nil {
		http.Error(w, "Erreur lors du traitement du like", http.StatusInternalServerError)
		return
	}

	// Récupère le nombre total de likes/dislikes
	likes, dislikes, err := models.GetCommentLikes(input.CommentID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des likes", http.StatusInternalServerError)
		return
	}

	// Crée une notification si l'auteur n'est pas l'utilisateur actuel
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

	// Retourne la réponse avec les nouveaux totaux
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"likes":    likes,
		"dislikes": dislikes,
	})
}
