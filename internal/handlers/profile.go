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

// Gestion de la page profil (Affichage et Modification)
func Profile(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	username = strings.Trim(username, "/")

	user, err := models.GetUserByUsername(username)
	if err != nil || user == nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	var currentUserID string
	if ctxUserID := r.Context().Value(security.ContextUserIDKey); ctxUserID != nil {
		currentUserID = ctxUserID.(string)
	}

	isOwner := currentUserID == user.ID

	var activities []models.Activity

	if r.Method == http.MethodGet {
		if isOwner {
			user.Warns, err = models.GetWarnsByUserID(user.ID)
			if err != nil {
				fmt.Println("Erreur de chargement des avertissements", err)
			}
		}

		if user.IsPublic || isOwner {
			activities, err = models.GetUserActivity(user.ID)
			if err != nil {
				fmt.Println("Erreur de chargement des activités", err)
			}
		}

		// Affichage du profil
		tmpl, err := template.ParseFiles("../public/template/profile.html")
		if err != nil {
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, struct {
			User       *models.User
			IsOwner    bool
			Activities []models.Activity
		}{
			User:       user,
			IsOwner:    isOwner,
			Activities: activities,
		})
		return
	}

	// Traitement de la modification (autorisé seulement pour le propriétaire)
	if r.Method == http.MethodPost && isOwner {
		var input struct {
			Username    string `json:"username"`
			Email       string `json:"email"`
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
			IsPublic    bool   `json:"is_public"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		// Vérifier si l'utilisateur s'est connecté avec Google
		if user.Provider.Valid && user.Provider.String == "google" || user.Provider.Valid && user.Provider.String == "github" {
			http.Error(w, "Modification des informations impossible pour les comptes Google ou Github", http.StatusForbidden)
			return
		}

		// Vérifier si le nouvel email est déjà utilisé
		existingUser, _ := models.GetUserByEmail(input.Email)
		if existingUser != nil && existingUser.ID != user.ID {
			http.Error(w, "Cet email est déjà utilisé", http.StatusBadRequest)
			return
		}

		// Vérifier si le nouvel username est déjà utilisé
		existingUser, _ = models.GetUserByUsername(input.Username)
		if existingUser != nil && existingUser.ID != user.ID {
			http.Error(w, "Ce nom d'utilisateur est déjà pris", http.StatusBadRequest)
			return
		}

		// Vérifier si l'utilisateur change son mot de passe
		var hashedPassword string
		if input.NewPassword != "" {
			// Vérifier si l'utilisateur a un mot de passe défini (pour éviter les NULL venant de Google)
			if !user.Password.Valid || user.Password.String == "" {
				http.Error(w, "Aucun mot de passe défini, modification impossible", http.StatusUnauthorized)
				return
			}

			// Vérifier l'ancien mot de passe
			err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(input.OldPassword))
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
			hashedPassword = user.Password.String // Garder l'ancien mot de passe si non modifié
		}

		// Mettre à jour le profil
		err = models.UpdateUserProfile(user.ID, input.Username, input.Email, hashedPassword, input.IsPublic)
		if err != nil {
			http.Error(w, "Erreur lors de la mise à jour du profil", http.StatusInternalServerError)
			return
		}

		// Répondre avec succès
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Profil mis à jour avec succès"))
	}
}
