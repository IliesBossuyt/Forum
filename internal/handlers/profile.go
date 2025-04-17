package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"Forum/internal/models"
	"Forum/internal/security"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Gère l'affichage et la modification du profil utilisateur
func Profile(w http.ResponseWriter, r *http.Request) {
	// Récupère le nom d'utilisateur depuis l'URL
	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	username = strings.Trim(username, "/")

	// Récupère les informations de l'utilisateur
	user, err := models.GetUserByUsername(username)
	if err != nil || user == nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	// Récupère l'ID de l'utilisateur actuel
	var currentUserID string
	if ctxUserID := r.Context().Value(security.ContextUserIDKey); ctxUserID != nil {
		currentUserID = ctxUserID.(string)
	}

	// Récupère le rôle de l'utilisateur actuel
	role := ""
	if r.Context().Value(security.ContextRoleKey) != nil {
		role = r.Context().Value(security.ContextRoleKey).(string)
	}

	isOwner := currentUserID == user.ID

	var activities []models.Activity

	// Gestion de la requête GET (affichage du profil)
	if r.Method == http.MethodGet {
		// Charge les avertissements si l'utilisateur est le propriétaire
		if isOwner {
			user.Warns, err = models.GetWarnsByUserID(user.ID)
			if err != nil {
				fmt.Println("Erreur de chargement des avertissements", err)
			}
		}

		// Charge les activités si le profil est public ou si l'utilisateur a les droits
		if user.IsPublic || isOwner || role == "admin" || role == "moderator" {
			activities, err = models.GetUserActivity(user.ID)
			if err != nil {
				fmt.Println("Erreur de chargement des activités", err)
			}
		}

		// Affiche le template du profil
		tmpl, err := template.ParseFiles("../public/template/profile.html")
		if err != nil {
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, struct {
			User        *models.User
			IsOwner     bool
			Activities  []models.Activity
			CurrentRole string
		}{
			User:        user,
			IsOwner:     isOwner,
			Activities:  activities,
			CurrentRole: role,
		})
		return
	}

	// Gestion de la requête POST (modification du profil)
	if r.Method == http.MethodPost && isOwner {
		// Structure pour décoder les données de modification
		var input struct {
			Username    string `json:"username"`
			Email       string `json:"email"`
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
			IsPublic    bool   `json:"is_public"`
		}

		// Décode les données de la requête
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		isExternal := user.Provider.Valid && (user.Provider.String == "google" || user.Provider.String == "github")

		// Cas : compte Google/GitHub
		if isExternal {
			// On vérifie s'il y a tentative de modifier autre chose que la visibilité
			if input.Username != user.Username || input.Email != user.Email || input.NewPassword != "" {
				http.Error(w, "Seule la visibilité du profil est modifiable pour les comptes Google/GitHub", http.StatusForbidden)
				return
			}

			// On met à jour uniquement la visibilité
			err := models.UpdateVisibilityOnly(user.ID, input.IsPublic)
			if err != nil {
				http.Error(w, "Erreur lors de la mise à jour de la visibilité", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Visibilité mise à jour avec succès"))
			return
		}

		// Vérifie si le nouvel email est déjà utilisé
		existingUser, _ := models.GetUserByEmail(input.Email)
		if existingUser != nil && existingUser.ID != user.ID {
			http.Error(w, "Cet email est déjà utilisé", http.StatusBadRequest)
			return
		}

		// Vérifie si le nouvel username est déjà utilisé
		existingUser, _ = models.GetUserByUsername(input.Username)
		if existingUser != nil && existingUser.ID != user.ID {
			http.Error(w, "Ce nom d'utilisateur est déjà pris", http.StatusBadRequest)
			return
		}

		// Gestion du changement de mot de passe
		var hashedPassword string
		if input.NewPassword != "" {
			// Vérifie si l'utilisateur a un mot de passe défini
			if !user.Password.Valid || user.Password.String == "" {
				http.Error(w, "Aucun mot de passe défini, modification impossible", http.StatusUnauthorized)
				return
			}

			// Vérifie l'ancien mot de passe
			err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(input.OldPassword))
			if err != nil {
				http.Error(w, "Ancien mot de passe incorrect", http.StatusUnauthorized)
				return
			}

			// Hash le nouveau mot de passe
			hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
				return
			}
			hashedPassword = string(hashedPasswordBytes)
		} else {
			hashedPassword = user.Password.String // Garde l'ancien mot de passe si non modifié
		}

		// Met à jour le profil
		err = models.UpdateUserProfile(user.ID, input.Username, input.Email, hashedPassword, input.IsPublic)
		if err != nil {
			http.Error(w, "Erreur lors de la mise à jour du profil", http.StatusInternalServerError)
			return
		}

		// Répond avec succès
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Profil mis à jour avec succès"))
	}
}
