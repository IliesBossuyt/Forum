package models

import "Forum/internal/database"

func ToggleCommentLike(userID string, commentID int, value int) (added bool, err error) {
	var existingValue int
	err = database.DB.QueryRow(
		"SELECT value FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID, commentID,
	).Scan(&existingValue)

	if err == nil {
		if existingValue == value {
			// Même valeur, on supprime (toggle off)
			_, err = database.DB.Exec("DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID, commentID)
			return false, err
		}

		// Valeur différente (ex: dislike -> like), on met à jour
		_, err = database.DB.Exec("UPDATE comment_likes SET value = ? WHERE user_id = ? AND comment_id = ?", value, userID, commentID)
		return value == 1, err
	}

	// Aucun like/dislike existant, insertion
	_, err = database.DB.Exec("INSERT INTO comment_likes (user_id, comment_id, value) VALUES (?, ?, ?)", userID, commentID, value)
	return value == 1, err
}

// Dans models/comment.go
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
