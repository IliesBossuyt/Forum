package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"Forum/internal/database"
)

// Structure d'un utilisateur
type User struct {
	ID       string         // Identifiant unique
	Username string         // Nom d'utilisateur
	Email    string         // Adresse email
	Password sql.NullString // Mot de passe (peut être NULL)
	Role     string         // Rôle de l'utilisateur
	GoogleID sql.NullString // ID Google (peut être NULL)
	GitHubID sql.NullString // ID GitHub
	Provider sql.NullString // Fournisseur d'authentification (peut être NULL)
	Banned   bool           // Statut de bannissement
	IsPublic bool           // Visibilité du profil
	Warns    []Warn         // Liste des avertissements
}

// Crée un nouvel utilisateur
func CreateUser(username, email, password string) error {
	// Génère un UUID unique
	id := uuid.New().String()

	// Hash le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insère l'utilisateur en base
	_, err = database.DB.Exec("INSERT INTO users (id, username, email, password, role) VALUES (?, ?, ?, ?, ?)", id, username, email, string(hashedPassword), "user")

	if err != nil {
		log.Println("Erreur lors de l'insertion de l'utilisateur :", err)
		return err
	}
	return nil
}

// Récupère un utilisateur par email
func GetUserByEmail(email string) (*User, error) {
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password, role, google_id, provider, github_id, banned, is_public FROM users WHERE email = ?", email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.GoogleID, &user.Provider, &user.GitHubID, &user.Banned, &user.IsPublic)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user.Normalize()

	return &user, nil
}

// Récupère un utilisateur par ID
func GetUserByID(userID string) (*User, error) {
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password, role, google_id, provider, github_id, banned, is_public FROM users WHERE id = ?", userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.GoogleID, &user.Provider, &user.GitHubID, &user.Banned, &user.IsPublic)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user.Normalize()

	return &user, nil
}

// Récupère un utilisateur par nom d'utilisateur
func GetUserByUsername(username string) (*User, error) {
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password, role, google_id, provider, github_id, banned, is_public FROM users WHERE username = ?", username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.GoogleID, &user.Provider, &user.GitHubID, &user.Banned, &user.IsPublic)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user.Normalize()

	return &user, nil
}

// Récupère un utilisateur par email ou nom d'utilisateur
func GetUserByIdentifier(identifier string) (*User, error) {
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password, role, google_id, provider, github_id, banned, is_public FROM users WHERE username = ? OR email = ?",
		identifier, identifier,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.GoogleID, &user.Provider, &user.GitHubID, &user.Banned, &user.IsPublic)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user.Normalize()

	return &user, nil
}

// Met à jour le profil d'un utilisateur
func UpdateUserProfile(userID, username, email, password string, isPublic bool) error {
	_, err := database.DB.Exec(`
		UPDATE users SET username = ?, email = ?, password = ?, is_public = ? WHERE id = ?
	`, username, email, password, isPublic, userID)
	return err
}

// Crée un utilisateur Google
func CreateGoogleUser(username, email, googleID string) error {
	// Génère un UUID unique
	id := uuid.New().String()

	_, err := database.DB.Exec(`
		INSERT INTO users (id, username, email, google_id, password, provider, role) 
		VALUES (?, ?, ?, ?, NULL, ?, ?)`, id, username, email, googleID, "google", "user")

	return err
}

// Crée un utilisateur GitHub
func CreateGitHubUser(username, email, githubID string) error {
	// Génère un UUID unique
	id := uuid.New().String()

	_, err := database.DB.Exec(`
		INSERT INTO users (id, username, email, github_id, password, provider, role) 
		VALUES (?, ?, ?, ?, NULL, ?, ?)`, id, username, email, githubID, "github", "user")

	return err
}

// Normalise les champs NULL de l'utilisateur
func (u *User) Normalize() {
	if !u.Password.Valid {
		u.Password.String = ""
	}
	if !u.GoogleID.Valid {
		u.GoogleID.String = ""
	}
	if !u.Provider.Valid {
		u.Provider.String = ""
	}
}

// Structure d'une activité utilisateur
type Activity struct {
	Type      string    // Type d'activité
	Content   string    // Contenu de l'activité
	Target    string    // Cible de l'activité
	CreatedAt time.Time // Date de création
}

// Récupère les activités d'un utilisateur
func GetUserActivity(userID string) ([]Activity, error) {
	// Requête pour obtenir toutes les activités
	rows, err := database.DB.Query(`
		SELECT 'post' AS type, content, NULL AS target, created_at
		FROM posts
		WHERE user_id = ?

		UNION ALL

		SELECT 'comment' AS type, content, 
			(SELECT content FROM posts WHERE posts.id = comments.post_id) AS target, created_at
		FROM comments
		WHERE author_id = ?

		UNION ALL

		SELECT 'like_post' AS type, '' AS content, 
			(SELECT content FROM posts WHERE posts.id = likes.post_id) AS target, created_at
		FROM likes
		WHERE user_id = ? AND value = 1

		UNION ALL

		SELECT 'dislike_post' AS type, '' AS content, 
			(SELECT content FROM posts WHERE posts.id = likes.post_id) AS target, created_at
		FROM likes
		WHERE user_id = ? AND value = -1

		UNION ALL

		SELECT 'like_comment' AS type, '' AS content, 
			(SELECT content FROM comments WHERE comments.id = comment_likes.comment_id) AS target, created_at
		FROM comment_likes
		WHERE user_id = ? AND value = 1

		UNION ALL

		SELECT 'dislike_comment' AS type, '' AS content, 
			(SELECT content FROM comments WHERE comments.id = comment_likes.comment_id) AS target, created_at
		FROM comment_likes
		WHERE user_id = ? AND value = -1

		ORDER BY created_at DESC
	`, userID, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourt et formate les résultats
	var activities []Activity
	for rows.Next() {
		var a Activity
		var rawTarget sql.NullString

		err := rows.Scan(&a.Type, &a.Content, &rawTarget, &a.CreatedAt)
		if err != nil {
			return nil, err
		}

		if rawTarget.Valid {
			a.Target = rawTarget.String
		} else {
			a.Target = "(contenu supprimé)"
		}

		activities = append(activities, a)
	}
	return activities, nil
}
