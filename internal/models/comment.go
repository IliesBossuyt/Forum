package models

import (
	"Forum/internal/database"
	"time"
)

// Structure d'un commentaire
type Comment struct {
	ID            int    // Identifiant unique
	PostID        int    // ID du post associé
	UserID        string // ID de l'auteur
	Username      string // Nom de l'auteur
	Content       string // Contenu du commentaire
	CreatedAt     string // Date de création
	Likes         int    // Nombre de likes
	Dislikes      int    // Nombre de dislikes
	CurrentUserID string // ID de l'utilisateur actuel
}

// Récupère les commentaires d'un post
func GetCommentsByPostID(postID int, currentUserID string) ([]Comment, error) {
	// Requête SQL pour récupérer les commentaires avec leurs likes/dislikes
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

	// Parcourt et formate les résultats
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

// Crée un nouveau commentaire
func CreateComment(postID int, userID, content string) (int64, time.Time, error) {
	// Insère le commentaire dans la base
	result, err := database.DB.Exec(`
		INSERT INTO comments (post_id, author_id, content, created_at)
		VALUES (?, ?, ?, ?)`,
		postID, userID, content, time.Now())
	if err != nil {
		return 0, time.Time{}, err
	}

	// Récupère l'ID du commentaire créé
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, time.Time{}, err
	}

	return lastID, time.Now(), nil
}

// Supprime un commentaire et ses données associées
func DeleteComment(commentID int) error {
	// Supprime les likes du commentaire
	if _, err := database.DB.Exec("DELETE FROM comment_likes WHERE comment_id = ?", commentID); err != nil {
		return err
	}

	// Supprime les signalements du commentaire
	if _, err := database.DB.Exec("DELETE FROM comment_reports WHERE comment_id = ?", commentID); err != nil {
		return err
	}

	// Supprime le commentaire
	_, err := database.DB.Exec("DELETE FROM comments WHERE id = ?", commentID)
	return err
}

// Récupère l'ID de l'auteur d'un commentaire
func GetCommentAuthorID(commentID int) (string, error) {
	var authorID string
	err := database.DB.QueryRow("SELECT author_id FROM comments WHERE id = ?", commentID).Scan(&authorID)
	return authorID, err
}

// Met à jour le contenu d'un commentaire
func UpdateCommentContent(commentID int, newContent string) error {
	_, err := database.DB.Exec("UPDATE comments SET content = ? WHERE id = ?", newContent, commentID)
	return err
}

// Récupère un commentaire par son ID
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
