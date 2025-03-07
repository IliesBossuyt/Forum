package handlers

import (
	"html/template"
	"net/http"

	"Forum/internal/models"
)

// Gestion de la page d'accueil (forum)
func Home(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Erreur de récupération des posts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("../public/template/home.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, posts)
}
