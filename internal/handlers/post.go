package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Gestion de la création d'un post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Vérifier si l'utilisateur est connecté (via le cookie)
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer le userID à partir du token de session
	userAgent := r.UserAgent()
	userID, valid := security.ValidateSecureToken(cookie.Value, userAgent)
	if !valid {
		security.DeleteCookie(w, cookie.Value)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Vérification du contenu du post
	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Le message ne peut pas être vide", http.StatusBadRequest)
		return
	}

	// Gestion de l'upload d'image
	var imageName string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageName = handler.Filename
		imagePath := filepath.Join("public/static/images", imageName)
		outFile, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Erreur lors de l'upload de l'image", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()
		_, err = outFile.ReadFrom(file)
	}

	// Insérer le post dans la base de données
	err = models.InsertPost(userID, content, imageName)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du post", http.StatusInternalServerError)
		return
	}

	// Rediriger vers /home après la publication
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
