package handlers

import (
	"Forum/internal/models"
	"encoding/json"
	"net/http"
)

// Récupère les avertissements d'un utilisateur
func GetUserWarns(w http.ResponseWriter, r *http.Request) {
	// Vérifie la méthode HTTP
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupère l'ID de l'utilisateur
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "ID utilisateur manquant", http.StatusBadRequest)
		return
	}

	// Récupère les avertissements
	warns, err := models.GetWarnsByUserID(userID)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Initialise un tableau vide si aucun avertissement
	if warns == nil {
		warns = []models.Warn{}
	}

	// Renvoie les avertissements en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warns)
}
