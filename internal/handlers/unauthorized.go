package handlers

import (
	"html/template"
	"net/http"
)

func UnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../public/template/unauthorized.html")
	if err != nil {
		http.Redirect(w, r, "/auth/unauthorized", http.StatusSeeOther)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	tmpl.Execute(w, nil)
}
