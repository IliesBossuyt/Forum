package models

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"Forum/internal/database"
)

// Structure User
type User struct {
	ID       string
	Username string
	Email    string
	Password sql.NullString // Peut être NULL
	Role     string
	GoogleID sql.NullString // Peut être NULL
	GitHubID sql.NullString
	Provider sql.NullString // Peut être NULL
	Banned   bool
}

// Fonction pour créer un utilisateur
func CreateUser(username, email, password string) error {
	// Générer un UUID
	id := uuid.New().String()

	// Hasher le mot de passe avec bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insérer l'utilisateur en base
	_, err = database.DB.Exec("INSERT INTO users (id, username, email, password, role) VALUES (?, ?, ?, ?, ?)", id, username, email, string(hashedPassword), "user")

	if err != nil {
		log.Println("Erreur lors de l'insertion de l'utilisateur :", err)
		return err
	}

	log.Println("Utilisateur ajouté :", username)
	return nil
}

// Trouver un utilisateur par email
func GetUserByEmail(email string) (*User, error) {
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password, role, google_id, provider, github_id, banned FROM users WHERE email = ?", email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.GoogleID, &user.Provider, &user.GitHubID, &user.Banned)

	if err == sql.ErrNoRows {
		return nil, nil // Aucun utilisateur trouvé
	} else if err != nil {
		return nil, err
	}

	user.Normalize()

	return &user, nil
}

// Récupérer un utilisateur par ID
func GetUserByID(userID string) (*User, error) {
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password, role, google_id, provider, github_id, banned FROM users WHERE id = ?", userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.GoogleID, &user.Provider, &user.GitHubID, &user.Banned)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user.Normalize()

	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password, role, google_id, provider, github_id, banned FROM users WHERE username = ?", username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.GoogleID, &user.Provider, &user.GitHubID, &user.Banned)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user.Normalize()

	return &user, nil
}

func GetUserByIdentifier(identifier string) (*User, error) {
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password, role, google_id, provider, github_id, banned FROM users WHERE username = ? OR email = ?",
		identifier, identifier,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.GoogleID, &user.Provider, &user.GitHubID, &user.Banned)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user.Normalize()

	return &user, nil
}

// Modifier le profil utilisateur
func UpdateUserProfile(userID, username, email, password string) error {
	_, err := database.DB.Exec("UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?",
		username, email, password, userID)
	return err
}

func CreateGoogleUser(username, email, googleID string) error {
	// Générer un UUID pour l'utilisateur
	id := uuid.New().String()

	_, err := database.DB.Exec(`
		INSERT INTO users (id, username, email, google_id, password, provider, role) 
		VALUES (?, ?, ?, ?, NULL, ?, ?)`, id, username, email, googleID, "google", "user")

	return err
}

func CreateGitHubUser(username, email, githubID string) error {
	// Générer un UUID pour l'utilisateur
	id := uuid.New().String()

	_, err := database.DB.Exec(`
		INSERT INTO users (id, username, email, github_id, password, provider, role) 
		VALUES (?, ?, ?, ?, NULL, ?, ?)`, id, username, email, githubID, "github", "user")

	return err
}

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
