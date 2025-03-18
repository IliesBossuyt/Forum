package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"

	"golang.org/x/crypto/bcrypt"
)

// 🔹 Gestion de la page profil (Affichage et Modification)
func Profile(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userAgent := r.UserAgent()
	userID, _, valid := security.ValidateSecureToken(cookie.Value, userAgent)
	if !valid {
		security.DeleteCookie(w, cookie.Value)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		// Charger et afficher le template
		tmpl, err := template.ParseFiles("../public/template/profile.html")
		if err != nil {
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, user)
		return
	}

	if r.Method == http.MethodPost {
		var input struct {
			Username    string `json:"username"`
			Email       string `json:"email"`
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		// Vérifier si le nouvel email est déjà utilisé
		existingUser, _ := models.GetUserByEmail(input.Email)
		if existingUser != nil && existingUser.ID != userID {
			http.Error(w, "Cet email est déjà utilisé", http.StatusBadRequest)
			return
		}

		// Vérifier si le nouvel username est déjà utilisé
		existingUser, _ = models.GetUserByUsername(input.Username)
		if existingUser != nil && existingUser.ID != userID {
			http.Error(w, "Ce nom d'utilisateur est déjà pris", http.StatusBadRequest)
			return
		}

		// Vérifier si l'utilisateur change son mot de passe
		var hashedPassword string
		if input.NewPassword != "" {
			// Vérifier l'ancien mot de passe
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword))
			if err != nil {
				http.Error(w, "Ancien mot de passe incorrect", http.StatusUnauthorized)
				return
			}

			// Hasher le nouveau mot de passe
			hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
				return
			}
			hashedPassword = string(hashedPasswordBytes)
		} else {
			hashedPassword = user.Password // Garder l'ancien mot de passe si non modifié
		}

		// Mettre à jour le profil
		err = models.UpdateUserProfile(userID, input.Username, input.Email, hashedPassword)
		if err != nil {
			http.Error(w, "Erreur lors de la mise à jour du profil", http.StatusInternalServerError)
			return
		}

		// Répondre avec succès
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Profil mis à jour avec succès"))
	}
}
