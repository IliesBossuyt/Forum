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
	Username  string // ✅ Ajouté
	Content   string
	CreatedAt time.Time
}

func InsertComment(userID string, postID int, content string) error {
	db := database.DB
	_, err := db.Exec(`
		INSERT INTO comments (post_id, user_id, content)
		VALUES (?, ?, ?)
	`, postID, userID, content)
	return err
}

func GetCommentsByPostID(postID int) ([]Comment, error) {
	db := database.DB

	rows, err := db.Query(`
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt); err != nil {
			log.Println("Erreur lecture commentaire:", err)
			continue
		}
		comments = append(comments, c)
	}
	return comments, nil
}
