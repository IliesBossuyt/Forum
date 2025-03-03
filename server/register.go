package server

import (
	"encoding/json"
	"net/http"
)

// Enregistrement (Register)
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "public/template/register.html")
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

		err = CreateUser(user.Username, user.Email, user.Password)
		if err != nil {
			http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("✅ Inscription réussie !"))
	}
}