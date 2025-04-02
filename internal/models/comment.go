package models

import (
	"Forum/internal/database"
	"log"
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    string
	Username  string
	Content   string
	CreatedAt time.Time
}

// Insérer un commentaire dans la base
func InsertComment(userID string, postID int, content string) error {
	_, err := database.DB.Exec(
		`INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`,
		postID, userID, content,
	)
	return err
}

// Récupérer les commentaires d'un post (avec nom d'utilisateur)
func GetCommentsByPostID(postID int) ([]Comment, error) {
	rows, err := database.DB.Query(`
		SELECT comments.id, comments.post_id, comments.user_id, users.username, comments.content, comments.created_at
		FROM comments
		JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
		ORDER BY comments.created_at DESC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt)
		if err != nil {
			log.Println("Erreur lecture commentaire:", err)
			continue
		}
		comments = append(comments, c)
	}
	return comments, nil
}