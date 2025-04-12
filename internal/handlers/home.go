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

	// Récupération des posts
	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Erreur de récupération des posts", http.StatusInternalServerError)
		return
	}

	// Associer les commentaires à chaque post
	for i := range posts {
		comments, err := models.GetCommentsByPostID(posts[i].ID, userID)
		if err != nil {
			http.Error(w, "Erreur de récupération des commentaires", http.StatusInternalServerError)
			return
		}
		posts[i].Comments = comments
		posts[i].CurrentUserID = userID
		posts[i].CurrentUserRole = role
	}

	// Struct pour passer au template
	data := struct {
		UserID string
		Role   string
		Posts  []models.Post
	}{
		UserID: userID,
		Role:   role,	
		Posts:  posts,
	}

	// Chargement du template HTML
	tmpl, err := template.ParseFiles("../public/template/home.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	// Exécution du template avec les données
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur d'exécution du template", http.StatusInternalServerError)
	}
}
