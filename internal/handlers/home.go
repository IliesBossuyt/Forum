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
	// Récupérer userID et rôle depuis le middleware
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)
	var username string
	if userID != "" {
		user, err := models.GetUserByID(userID)
		if err == nil && user != nil {
			username = user.Username
		}
	}	
	role, _ := r.Context().Value(security.ContextRoleKey).(string)

	// Récupération des catégories pour l'affichage
	categories, err := models.GetAllCategories()
	if err != nil {
		http.Error(w, "Erreur de récupération des catégories", http.StatusInternalServerError)
		return
	}

	// Vérifie s’il y a un filtre de catégorie dans l’URL
	categoryIDStr := r.URL.Query().Get("category")

	var posts []models.Post
	if categoryIDStr != "" {
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
		// Récupération de tous les posts
		posts, err = models.GetAllPosts()
		if err != nil {
			http.Error(w, "Erreur de récupération des posts", http.StatusInternalServerError)
			return
		}
	}

	// Associer les commentaires et infos utilisateur à chaque post
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

	// Struct pour le template
	data := struct {
		UserID     string
		Role       string
		Posts      []models.Post
		Categories []models.Category
		Username string
	}{
		UserID:     userID,
		Role:       role,
		Posts:      posts,
		Categories: categories,
		Username:  username,
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
