package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"Forum/internal/models"
)

// Gère l'inscription des utilisateurs
func Register(w http.ResponseWriter, r *http.Request) {
	// Gestion de la requête GET (affichage du formulaire)
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "../public/template/register.html")
		return
	} else if r.Method == http.MethodPost {
		// Structure pour décoder les données d'inscription
		var user struct {
			Username string `json:"username"` // Nom d'utilisateur
			Email    string `json:"email"`    // Adresse email
			Password string `json:"password"` // Mot de passe
		}

		// Décode les données de la requête
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		// Crée l'utilisateur dans la base de données
		err = models.CreateUser(user.Username, user.Email, user.Password)
		if err != nil {
			// Gestion des erreurs de contrainte d'unicité
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

		// Répond avec succès
		w.Write([]byte("Inscription réussie !"))
	}
}
