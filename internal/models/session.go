package models

import (
	"Forum/internal/database"
	"fmt"
	"time"
)

// Créer une session en base
func CreateSession(sessionUUID, userID, userAgent, role string, expiresAt time.Time) error {
	_, err := database.DB.Exec(
		"INSERT INTO sessions (token, user_id, user_agent, role, expires_at) VALUES (?, ?, ?, ?, ?)",
		sessionUUID, userID, userAgent, role, expiresAt,
	)
	return err
}

// Supprimer une session en base
func DeleteSession(sessionUUID string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE token = ?", sessionUUID)
	return err
}

func GetUserIDFromSession(sessionUUID string) (string, string, time.Time, error) {
	var userID, role string
	var expiresAt time.Time

	err := database.DB.QueryRow(
		"SELECT user_id, role, expires_at FROM sessions WHERE token = ?",
		sessionUUID,
	).Scan(&userID, &role, &expiresAt)

	if err != nil {
		return "", "", time.Time{}, err
	}

	return userID, role, expiresAt, nil
}

func UpdateUserSessionRole(userID, newRole string) error {
	_, err := database.DB.Exec(`
		UPDATE sessions SET role = ? WHERE user_id = ?
	`, newRole, userID)
	return err
}

func CleanExpiredSessions() {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE expires_at < NOW()")
	if err != nil {
		fmt.Println("Erreur nettoyage des sessions expirées :", err)
	}
}
