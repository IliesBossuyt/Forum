package server

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	_, err = DB.Exec("INSERT INTO users (id, username, email, password, role) VALUES (?, ?, ?, ?, ?)",
		id, username, email, string(hashedPassword), "user")

	if err != nil {
		log.Println("❌ Erreur lors de l'insertion de l'utilisateur :", err)
		return err
	}

	log.Println("✅ Utilisateur ajouté :", username)
	return nil
}

// Trouver un utilisateur par email
func GetUserByEmail(email string) (*User, error) {
	row := DB.QueryRow("SELECT id, username, email, password, role FROM users WHERE email = ?", email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return nil, nil // Aucun utilisateur trouvé
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id string) (*User, error) {
	row := DB.QueryRow("SELECT id, username, email, role FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
