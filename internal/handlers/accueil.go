package handlers

import (
	"html/template"
	"net/http"
)

func Accueil(w http.ResponseWriter, r *http.Request) {
	// Correction du chemin du template
	tmpl := template.Must(template.ParseFiles("../public/template/accueil.html"))

	data := struct{}{}

	tmpl.Execute(w, data)

}