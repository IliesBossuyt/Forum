package models

import (
	"Forum/internal/database"
	"time"
)

// Structure d'un signalement de post
type Report struct {
	ID           int       // Identifiant unique
	PostID       int       // ID du post signalé
	Reporter     string    // Nom de l'utilisateur qui signale
	Reason       string    // Raison du signalement
	CreatedAt    time.Time // Date de création
	PostContent  string    // Contenu du post signalé
	PostAuthor   string    // Nom de l'auteur du post
	PostAuthorID string    // ID de l'auteur du post
	PostImage    *string   // Image du post (pointeur pour accepter NULL)
}

// Crée un nouveau signalement de post
func CreateReport(postID int, reporterID string, reason string) error {
	_, err := database.DB.Exec(`
		INSERT INTO reports (post_id, reporter_id, reason, created_at)
		VALUES (?, ?, ?, ?)`,
		postID, reporterID, reason, time.Now(),
	)
	return err
}

// Récupère tous les signalements de posts
func GetAllReports() ([]Report, error) {
	// Requête pour obtenir les signalements avec les détails des posts
	rows, err := database.DB.Query(`
		SELECT 
			r.id,
			r.post_id,
			u.username AS reporter_name,
			r.reason,
			r.created_at,
			p.content AS post_content,
			pu.username AS post_author,
			p.image
		FROM reports r
		JOIN users u ON r.reporter_id = u.id
		JOIN posts p ON r.post_id = p.id
		JOIN users pu ON p.user_id = pu.id
		ORDER BY r.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourt et formate les résultats
	var reports []Report
	for rows.Next() {
		var r Report
		err := rows.Scan(&r.ID, &r.PostID, &r.Reporter, &r.Reason, &r.CreatedAt, &r.PostContent, &r.PostAuthor, &r.PostImage)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}

// Supprime un signalement de post par son ID
func DeleteReportByID(reportID int) error {
	_, err := database.DB.Exec("DELETE FROM reports WHERE id = ?", reportID)
	return err
}
