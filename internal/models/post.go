package models

import (
	"Forum/internal/database"
)

type Post struct {
	ID        int
	UserID    string
	Username  string
	Content   string
	Image     string
	CreatedAt string
}

// Récupérer tous les posts
func GetAllPosts() ([]Post, error) {
	rows, err := database.DB.Query(`
		SELECT posts.id, posts.user_id, users.username, posts.content, posts.image, posts.created_at 
		FROM posts 
		JOIN users ON posts.user_id = users.id 
		ORDER BY posts.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Content, &post.Image, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func InsertPost(userID, content, image string) error {
	_, err := database.DB.Exec("INSERT INTO posts (user_id, content, image) VALUES (?, ?, ?)", userID, content, image)
	return err
}
