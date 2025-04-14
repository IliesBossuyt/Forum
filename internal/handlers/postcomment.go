package handlers

import (
	"Forum/internal/database"
	"Forum/internal/models"
	"Forum/internal/security"
	"database/sql"
	"encoding/json"
	"net/http"
)

func PostComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(security.ContextUserIDKey).(string)

	var input struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Content == "" {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Insertion et récupération de l'ID
	commentID, createdAt, err := models.InsertComment(input.PostID, userID, input.Content)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}

	// Récupération du username
	var username string
	err = database.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			username = "Utilisateur inconnu"
		} else {
			http.Error(w, "Erreur lors de la récupération du nom", http.StatusInternalServerError)
			return
		}
	}

	commentIDInt := int(commentID)
	post, err := models.GetPostByID(input.PostID)
	if err == nil && post.UserID != userID {
		_ = models.CreateNotification(models.Notification{
			RecipientID: post.UserID,
			SenderID:    userID,
			Type:        "comment",
			PostID:      &input.PostID,
			CommentID:   &commentIDInt,
		})
	}

	// Réponse JSON
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"comment": map[string]interface{}{
			"id":        commentID,
			"content":   input.Content,
			"username":  username,
			"createdAt": createdAt.Format("02/01/2006 15:04"),
			"canEdit":   true,
			"canDelete": true, // c’est l'auteur
		},
	})
}
