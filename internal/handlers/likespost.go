package handlers

import (
	"encoding/json"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Gère les likes et dislikes sur les posts
func LikePost(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupère l'ID de l'utilisateur
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)

	// Structure pour le like/dislike
	var input struct {
		PostID int `json:"post_id"`
		Value  int `json:"value"` // 1 pour like, -1 pour dislike
	}

	// Décode les données de la requête
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Applique le like/dislike et récupère l'état
	added, wasLike, err := models.ToggleLike(userID, input.PostID, input.Value)
	if err != nil {
		http.Error(w, "Erreur lors du like/dislike", http.StatusInternalServerError)
		return
	}

	// Récupère le nombre total de likes/dislikes
	likes, dislikes, err := models.GetPostLikes(input.PostID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des likes/dislikes", http.StatusInternalServerError)
		return
	}

	// Crée une notification si l'auteur n'est pas l'utilisateur actuel
	if added {
		post, err := models.GetPostByID(input.PostID)
		if err == nil && post.UserID != userID {
			notifType := "dislike_post"
			if wasLike {
				notifType = "like_post"
			}

			models.CreateNotification(models.Notification{
				RecipientID: post.UserID,
				SenderID:    userID,
				Type:        notifType,
				PostID:      &input.PostID,
			})
		}
	}

	// Retourne la réponse avec les nouveaux totaux
	response := map[string]interface{}{
		"success":  true,
		"likes":    likes,
		"dislikes": dislikes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
