package models

import (
	"Forum/internal/database"
	"time"
)

// Ajouter ou modifier un like/dislike sur les posts
func ToggleLike(userID string, postID int, value int) (added bool, wasLike bool, err error) {
	var existingValue int
	err = database.DB.QueryRow(
		"SELECT value FROM likes WHERE user_id = ? AND post_id = ?", userID, postID,
	).Scan(&existingValue)

	if err == nil {
		if existingValue == value {
			// Même vote → suppression
			_, err = database.DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
			return false, false, err
		}

		// Changement de vote
		_, err = database.DB.Exec("UPDATE likes SET value = ? WHERE user_id = ? AND post_id = ?", value, userID, postID)
		return true, value == 1, err
	}

	// Premier vote
	_, err = database.DB.Exec(
		"INSERT INTO likes (user_id, post_id, value, created_at) VALUES (?, ?, ?, ?)",
		userID, postID, value, time.Now(),
	)
	return true, value == 1, err
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
