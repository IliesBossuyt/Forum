package handlers

import (
	"encoding/json"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
)

// 🔹 Handler pour liker/disliker un post
func LikePost(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	var input struct {
		PostID int `json:"post_id"`
		Value  int `json:"value"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Appliquer le like/dislike
	err = models.ToggleLike(userID, input.PostID, input.Value)
	if err != nil {
		http.Error(w, "Erreur lors du like/dislike", http.StatusInternalServerError)
		return
	}

	// Récupérer les nouvelles valeurs des likes/dislikes
	likes, dislikes, err := models.GetPostLikes(input.PostID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des likes/dislikes", http.StatusInternalServerError)
		return
	}

	// Renvoyer la réponse JSON
	response := map[string]interface{}{
		"success":  true,
		"likes":    likes,
		"dislikes": dislikes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
