package handlers

import (
	
	"html/template"
	"log"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
)

type PostView struct {
	Post     models.Post
	Comments []models.Comment
}

// Gestion de la page d'accueil (forum)
func Home(w http.ResponseWriter, r *http.Request) {
	// Récupérer userID et rôle depuis le middleware
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)
	role, _ := r.Context().Value(security.ContextRoleKey).(string)

	log.Printf("👤 Connexion utilisateur : userID=%s, role=%s", userID, role)

	// Récupération des posts
	posts, err := models.GetAllPosts()
	if err != nil {
		log.Println("❌ Erreur GetAllPosts:", err)
		http.Error(w, "Erreur de récupération des posts", http.StatusInternalServerError)
		return
	}

	var viewData []PostView
	for i := range posts {
		posts[i].CurrentUserID = userID
		posts[i].CurrentUserRole = role

		comments, err := models.GetCommentsByPostID(posts[i].ID)
		if err != nil {
			log.Printf("⚠️ Erreur GetCommentsByPostID pour post %d : %v", posts[i].ID, err)
		}

		viewData = append(viewData, PostView{
			Post:     posts[i],
			Comments: comments,
		})
	}

	tmpl, err := template.ParseFiles("../public/template/home.html")
	if err != nil {
		log.Println("❌ Erreur de chargement du template :", err)
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, viewData)
	if err != nil {
		log.Println("❌ Erreur rendering template :", err)
	}
}
