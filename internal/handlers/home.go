package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Gère l'affichage de la page d'accueil du forum
func Home(w http.ResponseWriter, r *http.Request) {
	// Récupère les informations de l'utilisateur connecté
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)
	var username string
	if userID != "" {
		user, err := models.GetUserByID(userID)
		if err == nil && user != nil {
			username = user.Username
		}
	}
	role, _ := r.Context().Value(security.ContextRoleKey).(string)

	// Récupère toutes les catégories
	categories, err := models.GetAllCategories()
	if err != nil {
		http.Error(w, "Erreur de récupération des catégories", http.StatusInternalServerError)
		return
	}

	// Récupère les paramètres de tri et de filtre
	sort := r.URL.Query().Get("sort")
	categoryIDStr := r.URL.Query().Get("category")

	var posts []models.Post

	if sort == "top" {
		// Tri par likes
		posts, err = models.GetTopPosts()
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des posts populaires", http.StatusInternalServerError)
			return
		}
	} else if categoryIDStr != "" {
		// Récupère les posts d'une catégorie spécifique
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
		// Récupère tous les posts
		posts, err = models.GetAllPosts()
		if err != nil {
			http.Error(w, "Erreur de récupération des posts", http.StatusInternalServerError)
			return
		}
	}

	// Ajoute les commentaires et les informations utilisateur à chaque post
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

	// Prépare les données pour le template
	data := struct {
		UserID     string
		Role       string
		Posts      []models.Post
		Categories []models.Category
		Username   string
	}{
		UserID:     userID,
		Role:       role,
		Posts:      posts,
		Categories: categories,
		Username:   username,
	}

	// Charge et exécute le template
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
