package models

import (
	"Forum/internal/database"
	"time"
)

type Comment struct {
	ID            int
	PostID        int
	UserID        string
	Username      string
	Content       string
	CreatedAt     string
	Likes         int
	Dislikes      int
	CurrentUserID string
}

func GetCommentsByPostID(postID int, currentUserID string) ([]Comment, error) {
	rows, err := database.DB.Query(`
		SELECT c.id, c.post_id, c.author_id, u.username, c.content, c.created_at,
			COALESCE(SUM(CASE WHEN cl.user_id IS NOT NULL THEN 1 ELSE 0 END), 0) AS likes
		FROM comments c
		JOIN users u ON c.author_id = u.id
		LEFT JOIN comment_likes cl ON c.id = cl.comment_id
		WHERE c.post_id = ?
		GROUP BY c.id
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		var rawTime time.Time
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.Likes)
		if err != nil {
			return nil, err
		}
		c.CreatedAt = rawTime.Format("02/01/2006 15:04")
		comments = append(comments, c)
	}
	return comments, nil
}

func InsertComment(postID int, userID, content string) error {
	_, err := database.DB.Exec("INSERT INTO comments (post_id, author_id, content) VALUES (?, ?, ?)", postID, userID, content)
	return err
}

func DeleteComment(commentID int) error {
	_, err := database.DB.Exec("DELETE FROM comment_likes WHERE comment_id = ?", commentID)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec("DELETE FROM comments WHERE id = ?", commentID)
	return err
}
