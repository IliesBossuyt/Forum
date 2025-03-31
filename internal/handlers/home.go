package handlers

import (
	"html/template"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Gestion de la page d'accueil (forum)
func Home(w http.ResponseWriter, r *http.Request) {
	// Récupérer userID et rôle depuis le middleware
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)
	role, _ := r.Context().Value(security.ContextRoleKey).(string)

	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Erreur de récupération des posts", http.StatusInternalServerError)
		return
	}

	// Passer CurrentUserID à chaque post
	for i := range posts {
		posts[i].CurrentUserID = userID
		posts[i].CurrentUserRole = role
	}

	tmpl, err := template.ParseFiles("../public/template/home.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, posts)
}
