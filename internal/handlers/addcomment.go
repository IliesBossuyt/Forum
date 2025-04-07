package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"Forum/internal/models"
	"Forum/internal/security"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupération du userID depuis le middleware
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)

	// Parse JSON (pas Form)
	var input struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Content == "" {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Insertion
	err := models.InsertComment(userID, input.PostID, input.Content)
	if err != nil {
		http.Error(w, "Erreur serveur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}

	// Récupérer l'utilisateur pour afficher son pseudo
	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Réponse JSON
	response := map[string]interface{}{
		"success":    true,
		"username":   user.Username,
		"created_at": time.Now().Format("02/01/2006 15:04"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
