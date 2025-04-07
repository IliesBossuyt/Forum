package models

import (
	"Forum/internal/database"
)

// Ajouter ou modifier un like/dislike sur les posts
func ToggleLike(userID string, postID int, value int) error {
	// Vérifier si l'utilisateur a déjà liké ou disliké ce post
	var existingValue int
	err := database.DB.QueryRow(
		"SELECT value FROM likes WHERE user_id = ? AND post_id = ?", userID, postID,
	).Scan(&existingValue)

	if err == nil { // L'utilisateur a déjà liké/disliké ce post
		if existingValue == value {
			// Si la valeur est identique, supprimer le like/dislike (annulation)
			_, err = database.DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
			return err
		}

		// Si la valeur est différente, on met à jour
		_, err = database.DB.Exec("UPDATE likes SET value = ? WHERE user_id = ? AND post_id = ?", value, userID, postID)
		return err
	}

	// Si l'utilisateur n'a pas encore liké, insérer un nouveau like/dislike
	_, err = database.DB.Exec("INSERT INTO likes (user_id, post_id, value) VALUES (?, ?, ?)", userID, postID, value)
	return err
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
