package handlers

import (
	"html/template"
	"net/http"
)

// Gère l'affichage de la page 404
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Charge le template de la page 404
	tmpl, err := template.ParseFiles("../public/template/notfound.html")
	if err != nil {
		http.Error(w, "Erreur 404", http.StatusNotFound)
		return
	}

	// Définit le code de statut 404
	w.WriteHeader(http.StatusNotFound)

	// Affiche la page 404
	tmpl.Execute(w, nil)
}

// Ajoute une gestion 404 personnalisée à un ServeMux
func WithNotFoundFallback(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vérifie si une route correspond à la requête
		handler, pattern := mux.Handler(r)
		if pattern == "" {
			// Affiche la page 404 si aucune route ne correspond
			NotFoundHandler(w, r)
			return
		}
		// Exécute le handler de la route trouvée
		handler.ServeHTTP(w, r)
	})
}
