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
	Password string
	Role     string
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
	row := database.DB.QueryRow("SELECT id, username, email, password, role FROM users WHERE email = ?", email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return nil, nil // Aucun utilisateur trouvé
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

// 🔹 Récupérer un utilisateur par ID
func GetUserByID(userID string) (*User, error) {
	var user User
	err := database.DB.QueryRow("SELECT id, username, email FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	row := database.DB.QueryRow("SELECT id, username, email, password, role FROM users WHERE username = ?", username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByIdentifier(identifier string) (*User, error) {
	var user User
	err := database.DB.QueryRow("SELECT id, username, email, password FROM users WHERE username = ? OR email = ?", identifier, identifier).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
