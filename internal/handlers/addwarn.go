package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

// Ajoute un avertissement à un utilisateur
func AddWarn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupère l'ID de l'utilisateur à avertir et la raison
	var input struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Récupère l'ID de l'utilisateur qui émet l'avertissement
	issuedBy := r.Context().Value(security.ContextUserIDKey).(string)

	// Ajoute l'avertissement dans la base de données
	err := models.AddWarn(input.UserID, issuedBy, input.Reason)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du warn", http.StatusInternalServerError)
		return
	}

	// Crée une notification
	if input.UserID != issuedBy {
		_ = models.CreateNotification(models.Notification{
			RecipientID: input.UserID,
			SenderID:    issuedBy,
			Type:        "warn",
		})
	}

	// Retourne le nouveau nombre total d'avertissements
	warns, _ := models.GetWarnsByUserID(input.UserID)
	count := len(warns)

	json.NewEncoder(w).Encode(map[string]any{
		"success":   true,
		"new_count": count,
	})
}
