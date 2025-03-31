package handlers

import (
	"Forum/internal/models"
	"net/http"
	"strconv"
	"strings"
)

// Récupérer l'image depuis la base
func GetImage(w http.ResponseWriter, r *http.Request) {
	// Extraire l’ID depuis l’URL manuellement
	path := r.URL.Path // ex: /image/5 ou /entry/image/5
	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		http.Error(w, "ID manquant dans l'URL", http.StatusBadRequest)
		return
	}
	postIDStr := segments[len(segments)-1]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	imageData, err := models.GetPostImage(postID)
	if err != nil {
		http.Error(w, "Image introuvable", http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageData)
}
