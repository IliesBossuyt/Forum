package models

import (
	"Forum/internal/database"
	"time"
)

// Structure d'un avertissement
type Warn struct {
	ID        int       // Identifiant unique
	UserID    string    // ID de l'utilisateur averti
	IssuedBy  string    // ID de l'administrateur qui a émis l'avertissement
	Issuer    string    // Nom de l'administrateur qui a émis l'avertissement
	Reason    string    // Raison de l'avertissement
	CreatedAt time.Time // Date de création
}

// Ajoute un avertissement à un utilisateur
func AddWarn(userID string, issuedBy string, reason string) error {
	_, err := database.DB.Exec(`
        INSERT INTO warns (user_id, issued_by, reason, created_at)
        VALUES (?, ?, ?, ?)`, userID, issuedBy, reason, time.Now())
	return err
}

// Récupère tous les avertissements d'un utilisateur
func GetWarnsByUserID(userID string) ([]Warn, error) {
	// Requête pour obtenir les avertissements avec le nom de l'administrateur
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

	// Parcourt et formate les résultats
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

// Récupère tous les avertissements
func GetAllWarns() ([]Warn, error) {
	// Requête pour obtenir tous les avertissements avec le nom des administrateurs
	rows, err := database.DB.Query(`
        SELECT w.id, w.user_id, w.issued_by, u.username AS issuer, w.reason, w.created_at
        FROM warns w
        JOIN users u ON w.issued_by = u.id
        ORDER BY w.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourt et formate les résultats
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

// Supprime un avertissement par son ID
func DeleteWarnByID(warnID int) error {
	_, err := database.DB.Exec("DELETE FROM warns WHERE id = ?", warnID)
	return err
}
