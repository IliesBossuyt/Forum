package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Gestion de la page d'accueil (forum)
func Home(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)
	role, _ := r.Context().Value(security.ContextRoleKey).(string)

	categories, err := models.GetAllCategories()
	if err != nil {
		http.Error(w, "Erreur de récupération des catégories", http.StatusInternalServerError)
		return
	}

	sort := r.URL.Query().Get("sort")
	categoryIDStr := r.URL.Query().Get("category")

	var posts []models.Post

	if sort == "top" {
		// 🔥 Tri par likes
		posts, err = models.GetTopPosts()
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des posts populaires", http.StatusInternalServerError)
			return
		}
	} else if categoryIDStr != "" {
		// Filtrage par catégorie
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "Catégorie invalide", http.StatusBadRequest)
			return
		}
		posts, err = models.GetPostsByCategoryID(categoryID)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des posts par catégorie", http.StatusInternalServerError)
			return
		}
	} else {
		// Tous les posts
		posts, err = models.GetAllPosts()
		if err != nil {
			http.Error(w, "Erreur de récupération des posts", http.StatusInternalServerError)
			return
		}
	}

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

	data := struct {
		UserID     string
		Role       string
		Posts      []models.Post
		Categories []models.Category
	}{
		UserID:     userID,
		Role:       role,
		Posts:      posts,
		Categories: categories,
	}

	tmpl, err := template.ParseFiles("../public/template/home.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur d'exécution du template", http.StatusInternalServerError)
	}
}
