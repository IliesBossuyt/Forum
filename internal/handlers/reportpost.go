package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

// Gère le signalement d'un post
func ReportPost(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Structure pour décoder les données du signalement
	var data struct {
		PostID int    `json:"post_id"` // ID du post signalé
		Reason string `json:"reason"`  // Raison du signalement
	}

	// Décode les données de la requête
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Récupère l'ID du post et de l'utilisateur qui signale
	postID := data.PostID
	userID := r.Context().Value(security.ContextUserIDKey).(string)

	// Crée le signalement dans la base de données
	err := models.CreateReport(postID, userID, data.Reason)
	if err != nil {
		http.Error(w, "Erreur lors du signalement", http.StatusInternalServerError)
		return
	}

	// Définit le type de contenu et renvoie la réponse JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
