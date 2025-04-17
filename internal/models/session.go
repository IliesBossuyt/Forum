package models

import (
	"Forum/internal/database"
	"fmt"
	"time"
)

// Crée une nouvelle session utilisateur
func CreateSession(sessionUUID, userID, userAgent, role string, expiresAt time.Time) error {
	_, err := database.DB.Exec(
		"INSERT INTO sessions (token, user_id, user_agent, role, expires_at) VALUES (?, ?, ?, ?, ?)",
		sessionUUID, userID, userAgent, role, expiresAt,
	)
	return err
}

// Supprime une session existante
func DeleteSession(sessionUUID string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE token = ?", sessionUUID)
	return err
}

// Récupère les informations d'une session
func GetUserIDFromSession(sessionUUID string) (string, string, time.Time, error) {
	var userID, role string
	var expiresAt time.Time

	// Récupère l'ID utilisateur, le rôle et la date d'expiration
	err := database.DB.QueryRow(
		"SELECT user_id, role, expires_at FROM sessions WHERE token = ?",
		sessionUUID,
	).Scan(&userID, &role, &expiresAt)

	if err != nil {
		return "", "", time.Time{}, err
	}

	return userID, role, expiresAt, nil
}

// Met à jour le rôle d'un utilisateur dans toutes ses sessions
func UpdateUserSessionRole(userID, newRole string) error {
	_, err := database.DB.Exec(`
		UPDATE sessions SET role = ? WHERE user_id = ?
	`, newRole, userID)
	return err
}

// Nettoie les sessions expirées
func CleanExpiredSessions() {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE expires_at < NOW()")
	if err != nil {
		fmt.Println("Erreur nettoyage des sessions expirées :", err)
	}
}
