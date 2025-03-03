package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"Forum/internal/models"
)

// Gestion de la page profil
func Profile(w http.ResponseWriter, r *http.Request) {
	// Vérifier si le cookie de session existe
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("Utilisateur non connecté, redirection vers /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer l'utilisateur via l'ID stocké dans le cookie
	user, err := models.GetUserByID(cookie.Value)
	if err != nil || user == nil {
		fmt.Println("Utilisateur introuvable, suppression du cookie et redirection vers /login")
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Charger le template `profile.html`
	tmpl, err := template.ParseFiles("public/template/profile.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	// Exécuter le template avec les données de l'utilisateur
	tmpl.Execute(w, user)
}
