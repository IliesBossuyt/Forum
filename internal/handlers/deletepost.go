package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Handler pour supprimer un post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	// Récupérer userID et rôle depuis le middleware
	userID := r.Context().Value(security.ContextUserIDKey).(string)
	role := r.Context().Value(security.ContextRoleKey).(string)

	// Lire l'ID du post depuis le JSON reçu
	var requestData struct {
		PostID string `json:"post_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Convertir l'ID du post en entier
	postID, err := strconv.Atoi(requestData.PostID)
	if err != nil {
		http.Error(w, "ID du post invalide", http.StatusBadRequest)
		return
	}

	// Vérifier si l'utilisateur est l'auteur du post ou un admin
	post, err := models.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	if userID != post.UserID && role != "admin" {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	// Supprimer le post de la base de données
	err = models.DeletePost(postID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	// Retourner une réponse JSON confirmant la suppression
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post supprimé avec succès",
	})
}
