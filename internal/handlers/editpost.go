package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
)

//modifier un post
func EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userAgent := r.UserAgent()
	userID, valid := security.ValidateSecureToken(cookie.Value, userAgent)
	if !valid {
		security.DeleteCookie(w, cookie.Value)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Token invalide, redirection vers /login")
		return
	}

	var input struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Content == "" {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Vérifier si l'utilisateur est bien le créateur du post
	post, err := models.GetPostByID(input.PostID)
	if err != nil || post.UserID != userID {
		http.Error(w, "Non autorisé", http.StatusForbidden)
		return
	}

	// Modifier le post
	err = models.UpdatePostContent(input.PostID, input.Content)
	if err != nil {
		http.Error(w, "Erreur lors de la modification du post", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
