package handlers

import (
	"net/http"
	"strconv"
	"log"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Gestion de l'ajout de commentaire
func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer userID depuis le middleware
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)

	// Récupération des données du formulaire
	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")

	// Valider les champs
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || content == "" {
		http.Error(w, "Entrées invalides", http.StatusBadRequest)
		return
	}

	log.Printf("🟢 Ajouter commentaire: userID=%s, postID=%d, content=\"%s\"\n", userID, postID, content)

	err = models.InsertComment(userID, postID, content)
	if err != nil {
		log.Println("❌ Erreur insertion commentaire:", err)
		http.Error(w, "Erreur serveur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page d'accueil avec les posts/commentaires
	http.Redirect(w, r, "/entry/home", http.StatusSeeOther)
}