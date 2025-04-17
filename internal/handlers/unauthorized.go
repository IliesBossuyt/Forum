package handlers

import (
	"html/template"
	"net/http"
)

// Affiche la page d'erreur 401
func UnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	// Charge la page d'erreur
	tmpl, err := template.ParseFiles("../public/template/unauthorized.html")
	if err != nil {
		http.Redirect(w, r, "/auth/unauthorized", http.StatusSeeOther)
		return
	}

	// Définit le code d'erreur 401
	w.WriteHeader(http.StatusUnauthorized)

	// Affiche la page
	tmpl.Execute(w, nil)
}
