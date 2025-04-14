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
		SELECT 
			c.id, 
			c.post_id, 
			c.author_id, 
			u.username, 
			c.content, 
			c.created_at,
			COALESCE(SUM(CASE WHEN cl.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN cl.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
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

		err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Username,
			&c.Content, &rawTime,
			&c.Likes, &c.Dislikes,
		)
		if err != nil {
			return nil, err
		}

		c.CreatedAt = rawTime.Format("02/01/2006 15:04")
		comments = append(comments, c)
	}
	return comments, nil
}

func InsertComment(postID int, userID, content string) (int64, time.Time, error) {
	loc, _ := time.LoadLocation("Europe/Paris")
	createdAt := time.Now().In(loc)

	result, err := database.DB.Exec(`
		INSERT INTO comments (post_id, author_id, content, created_at)
		VALUES (?, ?, ?, ?)`,
		postID, userID, content, createdAt)
	if err != nil {
		return 0, time.Time{}, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, time.Time{}, err
	}

	return lastID, createdAt, nil
}

func DeleteComment(commentID int) error {
	// Supprimer les likes
	if _, err := database.DB.Exec("DELETE FROM comment_likes WHERE comment_id = ?", commentID); err != nil {
		return err
	}

	// Supprimer les signalements liés
	if _, err := database.DB.Exec("DELETE FROM comment_reports WHERE comment_id = ?", commentID); err != nil {
		return err
	}

	// Supprimer le commentaire lui-même
	_, err := database.DB.Exec("DELETE FROM comments WHERE id = ?", commentID)
	return err
}

func GetCommentAuthorID(commentID int) (string, error) {
	var authorID string
	err := database.DB.QueryRow("SELECT author_id FROM comments WHERE id = ?", commentID).Scan(&authorID)
	return authorID, err
}

func UpdateCommentContent(commentID int, newContent string) error {
	_, err := database.DB.Exec("UPDATE comments SET content = ? WHERE id = ?", newContent, commentID)
	return err
}

func GetCommentByID(commentID int) (Comment, error) {
	var c Comment
	err := database.DB.QueryRow(`
		SELECT c.id, c.post_id, c.author_id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.id = ?
	`, commentID).Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt)

	return c, err
}
