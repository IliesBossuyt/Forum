package models

import (
	"Forum/internal/database"
	"time"
)

type Warn struct {
	ID        int
	UserID    string
	IssuedBy  string
	Issuer    string
	Reason    string
	CreatedAt time.Time
}

func AddWarn(userID string, issuedBy string, reason string) error {
	_, err := database.DB.Exec(`
        INSERT INTO warns (user_id, issued_by, reason)
        VALUES (?, ?, ?)`, userID, issuedBy, reason)
	return err
}

func GetWarnsByUserID(userID string) ([]Warn, error) {
	rows, err := database.DB.Query(`
        SELECT w.id, w.user_id, w.issued_by, u.username AS issuer, w.reason, w.created_at
        FROM warns w
        JOIN users u ON w.issued_by = u.id
        WHERE w.user_id = ?
        ORDER BY w.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warns []Warn
	for rows.Next() {
		var w Warn
		if err := rows.Scan(&w.ID, &w.UserID, &w.IssuedBy, &w.Issuer, &w.Reason, &w.CreatedAt); err != nil {
			return nil, err
		}
		warns = append(warns, w)
	}
	return warns, nil
}

func GetAllWarns() ([]Warn, error) {
	rows, err := database.DB.Query(`
        SELECT w.id, w.user_id, w.issued_by, u.username AS issuer, w.reason, w.created_at
        FROM warns w
        JOIN users u ON w.issued_by = u.id
        ORDER BY w.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warns []Warn
	for rows.Next() {
		var w Warn
		if err := rows.Scan(&w.ID, &w.UserID, &w.IssuedBy, &w.Issuer, &w.Reason, &w.CreatedAt); err != nil {
			return nil, err
		}
		warns = append(warns, w)
	}
	return warns, nil
}

func DeleteWarnByID(warnID int) error {
	_, err := database.DB.Exec("DELETE FROM warns WHERE id = ?", warnID)
	return err
}
