	package handlers

	import (
		"Forum/internal/models"
		"Forum/internal/security"
		"encoding/json"
		"net/http"
	)

	func DeleteComment(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		userID, _ := r.Context().Value(security.ContextUserIDKey).(string)
		role, _ := r.Context().Value(security.ContextRoleKey).(string)
		var input struct {
			CommentID int `json:"comment_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		// Récupérer l'auteur du commentaire
		authorID, err := models.GetCommentAuthorID(input.CommentID)
		if err != nil {
			http.Error(w, "Erreur lors de la vérification", http.StatusInternalServerError)
			return
		}

		// Autoriser si l'utilisateur est l'auteur OU admin OU moderateur
		if userID != authorID && role != "admin" && role != "moderator" {
			http.Error(w, "Non autorisé", http.StatusForbidden)
			return
		}

		// Supprimer le commentaire
		err = models.DeleteComment(input.CommentID)
		if err != nil {
			http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	}
