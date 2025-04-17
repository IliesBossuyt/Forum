package models

import (
	"Forum/internal/database"
	"time"
)

// Gère l'ajout/suppression d'un like/dislike sur un commentaire
func ToggleCommentLike(userID string, commentID int, value int) (added bool, wasLike bool, err error) {
	// Vérifie si l'utilisateur a déjà liké/disliké
	var existingValue int
	err = database.DB.QueryRow(
		"SELECT value FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID, commentID,
	).Scan(&existingValue)

	if err == nil {
		// Si même valeur, supprime le like/dislike
		if existingValue == value {
			_, err = database.DB.Exec("DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID, commentID)
			return false, false, err
		}

		// Si valeur différente, met à jour
		_, err = database.DB.Exec("UPDATE comment_likes SET value = ? WHERE user_id = ? AND comment_id = ?", value, userID, commentID)
		return true, value == 1, err
	}

	// Ajoute un nouveau like/dislike
	_, err = database.DB.Exec(`
	INSERT INTO comment_likes (user_id, comment_id, value, created_at)
	VALUES (?, ?, ?, ?)`,
		userID, commentID, value, time.Now(),
	)
	return true, value == 1, err
}

// Récupère le nombre de likes/dislikes d'un commentaire
func GetCommentLikes(commentID int) (int, int, error) {
	var likes, dislikes int
	err := database.DB.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0)
		FROM comment_likes
		WHERE comment_id = ?
	`, commentID).Scan(&likes, &dislikes)

	return likes, dislikes, err
}
