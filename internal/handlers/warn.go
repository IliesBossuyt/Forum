package handlers

import (
	"Forum/internal/models"
	"encoding/json"
	"net/http"
)

func GetUserWarns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "ID utilisateur manquant", http.StatusBadRequest)
		return
	}

	warns, err := models.GetWarnsByUserID(userID)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// S'il n'y a aucun warn, renvoyer un tableau vide :
	if warns == nil {
		warns = []models.Warn{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warns)
}
