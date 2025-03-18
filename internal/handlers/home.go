package handlers

import (
	"html/template"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Gestion de la page d'accueil (forum)
func Home(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté (pas obligatoire)
	cookie, err := r.Cookie("session")
	var currentUserID, currentUserRole string

	if err == nil { // Si le cookie existe, on tente de récupérer l'ID utilisateur et le rôle
		userAgent := r.UserAgent()
		currentUserID, currentUserRole, _ = security.ValidateSecureToken(cookie.Value, userAgent)
	}

	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Erreur de récupération des posts", http.StatusInternalServerError)
		return
	}

	// Passer CurrentUserID à chaque post
	for i := range posts {
		posts[i].CurrentUserID = currentUserID
		posts[i].CurrentUserRole = currentUserRole
	}

	tmpl, err := template.ParseFiles("../public/template/home.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, posts)
}
