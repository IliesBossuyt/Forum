package handlers

import (
	"io"
	"net/http"
	"strconv"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Formats d'image autorisés pour les posts
var allowedFormats = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/jpg":  true,
}

// Taille maximale autorisée pour les images (20 MB)
const maxFileSize = 20 * 1024 * 1024 // 20MB

// Gère la création d'un nouveau post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupère l'ID de l'utilisateur depuis le contexte
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)

	// Récupère le contenu texte du post
	content := r.FormValue("content")

	// Récupère l'image si elle est présente
	file, _, err := r.FormFile("image")
	var imageData []byte

	if err == nil { // Si un fichier est bien envoyé
		defer file.Close()
		imageData, _ = io.ReadAll(file) // Lit le fichier en bytes
	}

	// Vérifie qu'il y a au moins du texte ou une image
	if content == "" && len(imageData) == 0 {
		http.Error(w, "Le message doit contenir du texte ou une image", http.StatusBadRequest)
		return
	}

	// Traitement de l'image si elle est présente
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// Vérifie la taille du fichier
		if header.Size > maxFileSize {
			http.Error(w, "Fichier trop volumineux (max 20MB)", http.StatusBadRequest)
			return
		}

		// Vérifie le format du fichier
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture du fichier", http.StatusInternalServerError)
			return
		}

		mimeType := http.DetectContentType(buffer)
		if !allowedFormats[mimeType] {
			http.Error(w, "Format d'image non autorisé", http.StatusBadRequest)
			return
		}

		// Reviens au début du fichier pour le lire complètement
		file.Seek(0, io.SeekStart)
		imageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture de l'image", http.StatusInternalServerError)
			return
		}
	}

	// Crée le post dans la base de données
	postID, err := models.CreatePost(userID, content, imageData)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du post", http.StatusInternalServerError)
		return
	}

	// Récupère et traite les catégories associées au post
	categoryIDsStr := r.Form["categories"]
	var categoryIDs []int
	for _, idStr := range categoryIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Catégorie invalide", http.StatusBadRequest)
			return
		}
		categoryIDs = append(categoryIDs, id)
	}

	// Lie le post aux catégories sélectionnées
	err = models.LinkPostToCategories(postID, categoryIDs)
	if err != nil {
		http.Error(w, "Erreur lors du lien post-catégories", http.StatusInternalServerError)
		return
	}

	// Redirige vers la page d'accueil après la création
	http.Redirect(w, r, "/entry/home", http.StatusSeeOther)
}
