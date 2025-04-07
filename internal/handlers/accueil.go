package handlers

import (
	"html/template"
	"net/http"
)

// On vérifie que l'URL est bien exactement "/"
func Accueil(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFoundHandler(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("../public/template/accueil.html"))
	tmpl.Execute(w, nil)
}
