package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"Forum/internal/models"
)

// Enregistrement (Register)
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "../public/template/register.html")
		return
	} else if r.Method == http.MethodPost {
		var user struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		err = models.CreateUser(user.Username, user.Email, user.Password)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
				http.Error(w, "L'email est déjà utilisé", http.StatusConflict)
				return
			} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
				http.Error(w, "Le nom d'utilisateur est déjà pris", http.StatusConflict)
				return
			}
			http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Inscription réussie !"))
	}
}
