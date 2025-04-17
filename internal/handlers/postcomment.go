package handlers

import (
	"Forum/internal/database"
	"Forum/internal/models"
	"Forum/internal/security"
	"database/sql"
	"encoding/json"
	"net/http"
)

// Gère la création d'un nouveau commentaire
func PostComment(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupère l'ID de l'utilisateur depuis le contexte
	userID := r.Context().Value(security.ContextUserIDKey).(string)

	// Structure pour décoder les données du commentaire
	var input struct {
		PostID  int    `json:"post_id"` // ID du post commenté
		Content string `json:"content"` // Contenu du commentaire
	}

	// Décode les données de la requête
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Content == "" {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Crée le commentaire et récupère son ID et sa date de création
	commentID, createdAt, err := models.CreateComment(input.PostID, userID, input.Content)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}

	// Récupère le nom d'utilisateur de l'auteur
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

	// Crée une notification pour l'auteur du post si ce n'est pas le commentateur
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

	// Renvoie la réponse JSON avec les détails du commentaire
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"comment": map[string]interface{}{
			"id":        commentID,
			"content":   input.Content,
			"username":  username,
			"createdAt": createdAt.Format("02/01/2006 15:04"),
			"canEdit":   true,
			"canDelete": true, // L'auteur peut modifier/supprimer son commentaire
		},
	})
}
