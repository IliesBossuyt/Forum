package models

import (
	"Forum/internal/database"
	"time"
)

// Structure d'un post
type Post struct {
	ID              int        // Identifiant unique
	UserID          string     // ID de l'auteur
	Username        string     // Nom de l'auteur
	Content         string     // Contenu du post
	Image           []byte     // Image du post
	CreatedAt       string     // Date de création formatée
	Likes           int        // Nombre de likes
	Dislikes        int        // Nombre de dislikes
	CurrentUserID   string     // ID de l'utilisateur actuel
	CurrentUserRole string     // Rôle de l'utilisateur actuel
	Comments        []Comment  // Liste des commentaires
	Categories      []Category // Liste des catégories
}

// Récupère tous les posts avec leurs likes/dislikes
func GetAllPosts() ([]Post, error) {
	// Requête pour obtenir les posts avec leurs statistiques
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

	// Parcourt et formate les résultats
	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt time.Time
		err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Content, &post.Image, &createdAt, &post.Likes, &post.Dislikes)
		if err != nil {
			return nil, err
		}

		// Formatage de la date
		post.CreatedAt = createdAt.Format("02/01/2006 15:04")

		// Récupère les catégories du post
		categories, err := GetCategoriesByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}

	return posts, nil
}

// Crée un nouveau post
func CreatePost(userID, content string, image []byte) (int, error) {
	result, err := database.DB.Exec(
		"INSERT INTO posts (user_id, content, image, created_at) VALUES (?, ?, ?, ?)",
		userID, content, image, time.Now(),
	)
	if err != nil {
		return 0, err
	}

	// Récupère l'ID du post créé
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

// Récupère un post par son ID
func GetPostByID(postID int) (*Post, error) {
	var post Post
	err := database.DB.QueryRow("SELECT id, user_id, content FROM posts WHERE id = ?", postID).
		Scan(&post.ID, &post.UserID, &post.Content)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Met à jour un post existant
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

// Récupère l'image d'un post
func GetPostImage(postID int) ([]byte, error) {
	var imageData []byte
	err := database.DB.QueryRow("SELECT image FROM posts WHERE id = ?", postID).Scan(&imageData)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

// Supprime un post et toutes ses dépendances
func DeletePost(postID int) error {
	// Supprime les likes sur le post
	if _, err := database.DB.Exec("DELETE FROM likes WHERE post_id = ?", postID); err != nil {
		return err
	}

	// Supprime les likes sur les commentaires
	if _, err := database.DB.Exec("DELETE FROM comment_likes WHERE comment_id IN (SELECT id FROM comments WHERE post_id = ?)", postID); err != nil {
		return err
	}

	// Supprime les signalements sur les commentaires
	if _, err := database.DB.Exec("DELETE FROM comment_reports WHERE comment_id IN (SELECT id FROM comments WHERE post_id = ?)", postID); err != nil {
		return err
	}

	// Supprime les commentaires
	if _, err := database.DB.Exec("DELETE FROM comments WHERE post_id = ?", postID); err != nil {
		return err
	}

	// Supprime les signalements du post
	if _, err := database.DB.Exec("DELETE FROM reports WHERE post_id = ?", postID); err != nil {
		return err
	}

	// Supprime le post
	if _, err := database.DB.Exec("DELETE FROM posts WHERE id = ?", postID); err != nil {
		return err
	}

	return nil
}

// Récupère les catégories d'un post
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

// Récupère les posts d'une catégorie
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
		var createdAt time.Time
		err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Content, &post.Image, &createdAt, &post.Likes, &post.Dislikes)
		if err != nil {
			return nil, err
		}

		// Formatage de la date
		post.CreatedAt = createdAt.Format("02/01/2006 15:04")

		// Ajoute les catégories du post
		categories, err := GetCategoriesByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	return posts, nil
}

// Associe un post à des catégories
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

// Récupère les posts les plus populaires
func GetTopPosts() ([]Post, error) {
	rows, err := database.DB.Query(`
		SELECT p.id, p.user_id, u.username, p.content, p.image, p.created_at,
		       COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
		       COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
		GROUP BY p.id, p.user_id, u.username, p.content, p.image, p.created_at
		ORDER BY likes DESC
		LIMIT 20
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt time.Time
		err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Content, &post.Image, &createdAt, &post.Likes, &post.Dislikes)
		if err != nil {
			return nil, err
		}

		// Formatage de la date
		post.CreatedAt = createdAt.Format("02/01/2006 15:04")

		// Ajoute les catégories du post
		categories, err := GetCategoriesByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
		posts = append(posts, post)
	}
	return posts, nil
}
