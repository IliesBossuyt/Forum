package models

import (
	"Forum/internal/database"
	"time"
)

type Post struct {
	ID              int
	UserID          string
	Username        string
	Content         string
	Image           []byte
	CreatedAt       string
	Likes           int
	Dislikes        int
	CurrentUserID   string
	CurrentUserRole string
	Comments        []Comment
	Categories      []Category
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
		GROUP BY posts.id, posts.user_id, users.username, posts.content, posts.image, posts.created_at
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
		categories, err := GetCategoriesByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}

	return posts, nil
}

func CreatePost(userID, content string, image []byte) (int, error) {
	result, err := database.DB.Exec(
		"INSERT INTO posts (user_id, content, image, created_at) VALUES (?, ?, ?, ?)",
		userID, content, image, time.Now(),
	)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func GetPostByID(postID int) (*Post, error) {
	var post Post
	err := database.DB.QueryRow("SELECT id, user_id, content FROM posts WHERE id = ?", postID).
		Scan(&post.ID, &post.UserID, &post.Content)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Modifier un post
func UpdatePost(postID int, content string, imageData []byte, deleteImage bool) error {
	if deleteImage {
		_, err := database.DB.Exec("UPDATE posts SET content = ?, image = NULL WHERE id = ?", content, postID)
		return err
	}

	if len(imageData) > 0 {
		_, err := database.DB.Exec("UPDATE posts SET content = ?, image = ? WHERE id = ?", content, imageData, postID)
		return err
	}

	_, err := database.DB.Exec("UPDATE posts SET content = ? WHERE id = ?", content, postID)
	return err
}

// Récupérer l'image d'un post
func GetPostImage(postID int) ([]byte, error) {
	var imageData []byte
	err := database.DB.QueryRow("SELECT image FROM posts WHERE id = ?", postID).Scan(&imageData)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

// Supprimer un post (et ses dépendances)
func DeletePost(postID int) error {
	// Supprimer les likes
	if _, err := database.DB.Exec("DELETE FROM likes WHERE post_id = ?", postID); err != nil {
		return err
	}

	// Supprimer les commentaires
	if _, err := database.DB.Exec("DELETE FROM comment_likes WHERE comment_id IN (SELECT id FROM comments WHERE post_id = ?)", postID); err != nil {
		return err
	}
	if _, err := database.DB.Exec("DELETE FROM comments WHERE post_id = ?", postID); err != nil {
		return err
	}

	// Supprimer les signalements
	if _, err := database.DB.Exec("DELETE FROM reports WHERE post_id = ?", postID); err != nil {
		return err
	}

	// Supprimer le post lui-même
	if _, err := database.DB.Exec("DELETE FROM posts WHERE id = ?", postID); err != nil {
		return err
	}

	return nil
}

func GetCategoriesByPostID(postID int) ([]Category, error) {
	rows, err := database.DB.Query(`
		SELECT c.id, c.name
		FROM category c
		INNER JOIN post_category pc ON c.id = pc.category_id
		WHERE pc.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}


func GetPostsByCategoryID(categoryID int) ([]Post, error) {
	rows, err := database.DB.Query(`
		SELECT p.id, p.user_id, u.username, p.content, p.image, p.created_at,
		       COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
		       COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
		INNER JOIN post_category pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		GROUP BY p.id, p.user_id, u.username, p.content, p.image, p.created_at
		ORDER BY p.created_at DESC
	`, categoryID)

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

		// Ajouter les catégories du post
		categories, err := GetCategoriesByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	return posts, nil
}

func LinkPostToCategories(postID int, categoryIDs []int) error {
	for _, categoryID := range categoryIDs {
		_, err := database.DB.Exec(`
			INSERT INTO post_category (post_id, category_id)
			VALUES (?, ?)`, postID, categoryID)
		if err != nil {
			return err
		}
	}
	return nil
}
