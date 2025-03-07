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
	Likes     int
	Dislikes  int
}

// Récupérer tous les posts
func GetAllPosts() ([]Post, error) {
	rows, err := database.DB.Query(`
		SELECT posts.id, posts.user_id, users.username, posts.content, posts.image, posts.created_at,
		       COALESCE(SUM(CASE WHEN likes.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
		       COALESCE(SUM(CASE WHEN likes.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM posts
		JOIN users ON posts.user_id = users.id
		LEFT JOIN likes ON posts.id = likes.post_id
		GROUP BY posts.id
		ORDER BY posts.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Content, &post.Image, &post.CreatedAt, &post.Likes, &post.Dislikes)
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
