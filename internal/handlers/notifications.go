package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
)

// Récupère et renvoie les notifications d'un utilisateur au format JSON
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur depuis le contexte
	userID := r.Context().Value(security.ContextUserIDKey).(string)

	// Récupère les notifications de l'utilisateur
	notifs, err := models.GetNotificationsByUser(userID)
	if err != nil {
		http.Error(w, "Erreur de récupération des notifications", http.StatusInternalServerError)
		return
	}

	// Structure pour le format JSON des notifications
	type JSONNotif struct {
		Message   string `json:"message"`    // Message de la notification
		CreatedAt string `json:"created_at"` // Date de création formatée
		Seen      bool   `json:"seen"`       // État de lecture
	}

	// Convertit les notifications en format JSON
	var result = make([]JSONNotif, 0)
	for _, n := range notifs {
		result = append(result, JSONNotif{
			Message:   n.Message,
			CreatedAt: n.CreatedAt.Format("02/01/2006 15:04"),
			Seen:      n.Seen,
		})
	}

	// Envoie la réponse en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Marque toutes les notifications d'un utilisateur comme lues
func MarkNotificationsRead(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur depuis le contexte
	userID := r.Context().Value(security.ContextUserIDKey).(string)

	// Met à jour l'état des notifications
	if err := models.MarkNotificationsAsSeen(userID); err != nil {
		http.Error(w, "Erreur lors de la mise à jour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Supprime toutes les notifications d'un utilisateur
func DeleteAllNotifications(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur depuis le contexte
	userID := r.Context().Value(security.ContextUserIDKey).(string)

	// Supprime toutes les notifications
	err := models.DeleteAllNotificationsForUser(userID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	// Renvoie un tableau vide en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]any{})
}
