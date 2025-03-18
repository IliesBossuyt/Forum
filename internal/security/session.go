package security

import (
	"Forum/internal/database"
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

func GetUserIDFromSession(sessionUUID string) (string, string, error) {
	var userID, role string
	err := database.DB.QueryRow("SELECT user_id, role FROM sessions WHERE token = ?", sessionUUID).Scan(&userID, &role)
	if err != nil {
		return "", "", err
	}
	return userID, role, nil
}
