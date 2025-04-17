package models

import (
	"Forum/internal/database"
	"time"
)

// Gère l'ajout/suppression d'un like/dislike sur un post
func ToggleLike(userID string, postID int, value int) (added bool, wasLike bool, err error) {
	// Vérifie si l'utilisateur a déjà liké/disliké
	var existingValue int
	err = database.DB.QueryRow(
		"SELECT value FROM likes WHERE user_id = ? AND post_id = ?", userID, postID,
	).Scan(&existingValue)

	if err == nil {
		// Si même valeur, supprime le like/dislike
		if existingValue == value {
			_, err = database.DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
			return false, false, err
		}

		// Si valeur différente, met à jour
		_, err = database.DB.Exec("UPDATE likes SET value = ? WHERE user_id = ? AND post_id = ?", value, userID, postID)
		return true, value == 1, err
	}

	// Ajoute un nouveau like/dislike
	_, err = database.DB.Exec(
		"INSERT INTO likes (user_id, post_id, value, created_at) VALUES (?, ?, ?, ?)",
		userID, postID, value, time.Now(),
	)
	return true, value == 1, err
}

// Récupère le nombre de likes/dislikes d'un post
func GetPostLikes(postID int) (int, int, error) {
	var likes, dislikes int
	err := database.DB.QueryRow(
		`SELECT 
            COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0) AS likes, 
            COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0) AS dislikes 
        FROM likes WHERE post_id = ?`, postID,
	).Scan(&likes, &dislikes)

	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}
