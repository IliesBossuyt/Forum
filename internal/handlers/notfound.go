package handlers

import (
	"html/template"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../public/template/notfound.html")
	if err != nil {
		http.Error(w, "Erreur 404", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, nil)
}

// WithNotFoundFallback permet d'ajouter une gestion 404 locale à un ServeMux
func WithNotFoundFallback(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, pattern := mux.Handler(r)
		if pattern == "" {
			NotFoundHandler(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
