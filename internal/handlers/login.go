package handlers

import (
	"encoding/json"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {

	// Vérifier si l'utilisateur est déjà connecté
	cookie, err := r.Cookie("session")
	if err == nil {
		// Vérifier la validité du token dans le cookie
		userID, _, valid := security.ValidateSecureToken(cookie.Value, r.UserAgent())
		if valid {
			// Vérification de l'existence de l'utilisateur en base
			user, err := models.GetUserByID(userID)
			if err == nil && user != nil {
				// Rediriger directement vers le profil
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			}
		}
	}

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "../public/template/login.html")
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

		// Vérifier que l'utilisateur a bien un mot de passe (pour éviter les NULL venant de Google)
		if !user.Password.Valid || user.Password.String == "" {
			http.Error(w, "Aucun mot de passe défini pour cet utilisateur", http.StatusUnauthorized)
			return
		}

		// Vérifier le mot de passe
		err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(input.Password))
		if err != nil {
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Stocker l'ID et le rôle dans le cookie
		err = security.CreateCookie(w, r, user.ID, user.Role)
		if err != nil {
			http.Error(w, "Erreur lors de la création du cookie", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Connexion réussie !"))
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		security.DeleteCookie(w, cookie.Value)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
