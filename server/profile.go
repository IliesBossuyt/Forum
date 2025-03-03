package server

import (
	"html/template"
	"net/http"
)

// Gestion de la page profil
func Profile(w http.ResponseWriter, r *http.Request) {
	// Vérifier si le cookie de session existe
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer l'utilisateur via l'ID stocké dans le cookie
	user, err := GetUserByID(cookie.Value)
	if err != nil || user == nil {
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
