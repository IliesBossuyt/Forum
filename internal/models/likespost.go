package models

import (
	"Forum/internal/database"
)

// Ajouter ou modifier un like/dislike sur les posts
func ToggleLike(userID string, postID int, value int) (bool, error) {
	var existingValue int
	err := database.DB.QueryRow(
		"SELECT value FROM likes WHERE user_id = ? AND post_id = ?", userID, postID,
	).Scan(&existingValue)

	if err == nil {
		if existingValue == value {
			// L'utilisateur retire son vote
			_, err = database.DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
			return false, err
		}

		// L'utilisateur change son vote (like ↔ dislike)
		_, err = database.DB.Exec("UPDATE likes SET value = ? WHERE user_id = ? AND post_id = ?", value, userID, postID)
		return value == 1, err // seulement true si c'est un like
	}

	// Nouveau vote
	_, err = database.DB.Exec("INSERT INTO likes (user_id, post_id, value) VALUES (?, ?, ?)", userID, postID, value)
	return value == 1, err
}


// Récupérer le nombre de likes et dislikes d'un post
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
