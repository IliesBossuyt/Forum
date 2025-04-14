package handlers

import (
	"Forum/internal/models"
	"Forum/internal/security"
	"encoding/json"
	"net/http"
	"time"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(security.ContextUserIDKey).(string)

	notifs, err := models.GetNotificationsByUser(userID)
	if err != nil {
		http.Error(w, "Erreur de récupération des notifications", http.StatusInternalServerError)
		return
	}

	loc, _ := time.LoadLocation("Europe/Paris")
	type JSONNotif struct {
		Message   string `json:"message"`
		CreatedAt string `json:"created_at"`
		Seen      bool   `json:"seen"`
	}

	var result []JSONNotif
	for _, n := range notifs {
		result = append(result, JSONNotif{
			Message:   n.Message,
			CreatedAt: n.CreatedAt.In(loc).Format("02/01/2006 15:04"),
			Seen:      n.Seen,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func MarkNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(security.ContextUserIDKey).(string)
	if err := models.MarkNotificationsAsSeen(userID); err != nil {
		http.Error(w, "Erreur lors de la mise à jour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
