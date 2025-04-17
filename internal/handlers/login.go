package handlers

import (
	"encoding/json"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"

	"golang.org/x/crypto/bcrypt"
)

// Gère la connexion des utilisateurs
func Login(w http.ResponseWriter, r *http.Request) {
	// Gestion de la requête GET (affichage de la page de connexion)
	if r.Method == http.MethodGet {
		// Vérifie si l'utilisateur est déjà connecté
		cookie, err := r.Cookie("session")
		if err == nil {
			userID, _, valid := security.ValidateSecureToken(cookie.Value, r.UserAgent())
			if valid {
				user, err := models.GetUserByID(userID)
				if err == nil && user != nil && !user.Banned {
					http.Redirect(w, r, "/entry/home", http.StatusSeeOther)
					return
				}
			}

			// Nettoie le cookie si invalide ou utilisateur banni
			security.DeleteCookie(w, cookie.Value)
		}

		// Affiche la page de connexion
		http.ServeFile(w, r, "../public/template/login.html")
		return
	} else if r.Method == http.MethodPost {
		// Gestion de la requête POST (traitement de la connexion)
		var input struct {
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}

		// Décode les données de connexion
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		// Récupère l'utilisateur par email ou username
		user, err := models.GetUserByIdentifier(input.Identifier)
		if err != nil || user == nil {
			http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
			return
		}

		// Vérifie si l'utilisateur est banni
		if user.Banned {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"banned": "true",
			})
			return
		}

		// Vérifie la présence d'un mot de passe
		if !user.Password.Valid || user.Password.String == "" {
			http.Error(w, "Aucun mot de passe défini pour cet utilisateur", http.StatusUnauthorized)
			return
		}

		// Vérifie le mot de passe
		err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(input.Password))
		if err != nil {
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Crée le cookie de session
		err = security.CreateCookie(w, r, user.ID, user.Role)
		if err != nil {
			http.Error(w, "Erreur lors de la création du cookie", http.StatusInternalServerError)
			return
		}

		// Retourne la confirmation de connexion
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "Connexion réussie",
			"username": user.Username,
		})
	}
}

// Gère la déconnexion des utilisateurs
func Logout(w http.ResponseWriter, r *http.Request) {
	// Supprime le cookie de session s'il existe
	cookie, err := r.Cookie("session")
	if err == nil {
		security.DeleteCookie(w, cookie.Value)
	}

	// Redirige vers la page de connexion
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
