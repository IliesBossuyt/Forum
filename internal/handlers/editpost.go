package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Modifier un post (contenu + image)
func EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Vérifier l'authentification
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userAgent := r.UserAgent()
	userID, _, valid := security.ValidateSecureToken(cookie.Value, userAgent)
	if !valid {
		security.DeleteCookie(w, cookie.Value)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Token invalide, redirection vers /login")
		return
	}

	// Récupérer les valeurs du formulaire
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	deleteImage := r.FormValue("delete_image") == "true"

	// Vérifier si l'utilisateur est bien l'auteur du post
	post, err := models.GetPostByID(postID)
	if err != nil || post.UserID != userID {
		http.Error(w, "Non autorisé", http.StatusForbidden)
		return
	}

	// Gestion de l'image
	var imageData []byte
	imageUpdated := false
	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture de l'image", http.StatusInternalServerError)
			return
		}
		imageUpdated = true
	}

	// Mettre à jour le post (contenu + image)
	err = models.UpdatePost(postID, content, imageData, deleteImage)
	if err != nil {
		http.Error(w, "Erreur lors de la modification du post", http.StatusInternalServerError)
		return
	}

	// Répondre en JSON pour mise à jour en front
	response := map[string]interface{}{
		"success":      true,
		"imageUpdated": imageUpdated,
		"imageDeleted": deleteImage,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
