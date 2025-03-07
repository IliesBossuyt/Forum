package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
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

	// Vérifier le token avec User-Agent et IP actuels
	userAgent := r.UserAgent()
	userIP := security.ExtractIP(r.RemoteAddr)


	userID, valid := security.ValidateSecureToken(cookie.Value, userAgent, userIP)
	if !valid {
		fmt.Println("Session suspecte détectée ! Suppression du cookie et redirection.")
		security.DeleteCookie(w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer les infos utilisateur
	user, err := models.GetUserByID(userID)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Charger le template profile.html
	tmpl, err := template.ParseFiles("../public/template/profile.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	// Exécuter le template avec les données de l'utilisateur
	tmpl.Execute(w, user)
}
