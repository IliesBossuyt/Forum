package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

func AddWarn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	issuedBy := r.Context().Value(security.ContextUserIDKey).(string)

	err := models.AddWarn(input.UserID, issuedBy, input.Reason)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du warn", http.StatusInternalServerError)
		return
	}

	warns, _ := models.GetWarnsByUserID(input.UserID)
	count := len(warns)
	
	json.NewEncoder(w).Encode(map[string]any{
		"success":   true,
		"new_count": count, 
	})	
}

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
