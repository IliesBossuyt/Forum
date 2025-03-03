package handlers

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	// Supprimer le cookie de session	
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Rediriger vers la page de connexion
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
