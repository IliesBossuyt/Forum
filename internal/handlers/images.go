package handlers

import (
	"Forum/internal/models"
	"net/http"
	"strconv"
)

// Récupérer l'image depuis la base
func GetImage(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	imageData, err := models.GetPostImage(postID)
	if err != nil {
		http.Error(w, "Image introuvable", http.StatusNotFound)
		return
	}

	// Empêcher le cache du navigateur
	w.Header().Set("Cache-Control", "no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.Header().Set("Content-Type", "image/jpeg") // Ajuste selon le type réel
	w.Write(imageData)
}
