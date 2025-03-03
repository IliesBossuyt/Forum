package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Connexion (Login)
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "public/template/login.html")
		return
	} else if r.Method == http.MethodPost {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		user, err := GetUserByEmail(input.Email)
		if err != nil || user == nil {
			http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
			return
		}

		// Vérification du mot de passe
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Créer un cookie de session
		sessionToken := uuid.New().String()
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		w.Write([]byte("✅ Connexion réussie !"))
	}
}
