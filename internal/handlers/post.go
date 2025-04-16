package handlers

import (
	"io"
	"net/http"

	"Forum/internal/models"
	"Forum/internal/security"
)

// Formats autorisés
var allowedFormats = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/jpg":  true,
}

// Taille max des fichiers (20 MB)
const maxFileSize = 20 * 1024 * 1024 // 20MB

// Gestion de la création d'un post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer userID et rôle depuis le middleware
	userID, _ := r.Context().Value(security.ContextUserIDKey).(string)
	
	// Vérification du contenu du post : il faut AU MOINS du texte ou une image
	content := r.FormValue("content")

	// Récupération de l'image envoyée
	file, _, err := r.FormFile("image")
	var imageData []byte

	if err == nil { // Si un fichier est bien envoyé
		defer file.Close()
		imageData, _ = io.ReadAll(file) // Lire le fichier en bytes
	}

	// Vérification : soit `content`, soit `imageData`, soit les deux doivent être présents
	if content == "" && len(imageData) == 0 {
		http.Error(w, "Le message doit contenir du texte ou une image", http.StatusBadRequest)
		return
	}

	// Récupérer l’image et la stocker en BLOB avec vérifications
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// Vérifier la taille
		if header.Size > maxFileSize {
			http.Error(w, "Fichier trop volumineux (max 20MB)", http.StatusBadRequest)
			return
		}

		// Vérifier le format du fichier
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

		// Revenir au début du fichier avant de le lire complètement
		file.Seek(0, io.SeekStart)
		imageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture de l'image", http.StatusInternalServerError)
			return
		}
	}

	// Insérer le post dans la base de données
	err = models.CreatePost(userID, content, imageData)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du post", http.StatusInternalServerError)
		return
	}
	// Rediriger vers /home après la publication
	http.Redirect(w, r, "/entry/home", http.StatusSeeOther)
}
