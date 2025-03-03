package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"Forum/internal/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "public/template/login.html")
		return
	} else if r.Method == http.MethodPost {
		var input struct {
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		// Récupérer l'utilisateur soit par email, soit par username
		user, err := models.GetUserByIdentifier(input.Identifier)
		if err != nil || user == nil {
			http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
			return
		}

		// Vérifier le mot de passe
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Créer un cookie de session avec l'ID utilisateur
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   user.ID, // On stocke directement l'ID utilisateur
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		w.Write([]byte("Connexion réussie !"))
	}
}
