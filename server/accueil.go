package server

import (
	"html/template"
	"net/http"
)

func (forum *Forum) Accueil(w http.ResponseWriter, r *http.Request) {
	// Correction du chemin du template
	tmpl := template.Must(template.ParseFiles("public/template/accueil.html"))

	data := Forum{}

	tmpl.Execute(w, data)

}