package models

import "Forum/internal/database"

// Like ou Dislike un commentaire (ou annule si déjà fait)
func ToggleCommentLike(userID string, commentID int, value int) error {
	var existingValue int
	err := database.DB.QueryRow(
		"SELECT value FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID, commentID,
	).Scan(&existingValue)

	if err == nil {
		if existingValue == value {
			_, err = database.DB.Exec("DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID, commentID)
			return err
		}
		_, err = database.DB.Exec("UPDATE comment_likes SET value = ? WHERE user_id = ? AND comment_id = ?", value, userID, commentID)
		return err
	}

	_, err = database.DB.Exec("INSERT INTO comment_likes (user_id, comment_id, value) VALUES (?, ?, ?)", userID, commentID, value)
	return err
}

// Récupère le total de likes / dislikes d’un commentaire
func GetCommentLikes(commentID int) (int, int, error) {
	var likes, dislikes int
	err := database.DB.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comment_likes WHERE comment_id = ?
	`, commentID).Scan(&likes, &dislikes)

	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}
