package models

import (
	"Forum/internal/database"
	"time"
)

// Structure d'un signalement de commentaire
type CommentReport struct {
	ID              int       // Identifiant unique
	CommentID       int       // ID du commentaire signalé
	Reporter        string    // Nom de l'utilisateur qui signale
	Reason          string    // Raison du signalement
	CreatedAt       time.Time // Date de création
	CommentText     string    // Contenu du commentaire signalé
	CommentAuthor   string    // Nom de l'auteur du commentaire
	CommentAuthorID string    // ID de l'auteur du commentaire
}

// Récupère tous les signalements de commentaires
func GetAllCommentReports() ([]CommentReport, error) {
	// Requête pour obtenir les signalements avec les détails des commentaires
	rows, err := database.DB.Query(`
		SELECT cr.id, cr.comment_id, u.username AS reporter_name, cr.reason, cr.created_at,
			c.content AS comment_content, cu.username AS comment_author, cu.id AS comment_author_id
		FROM comment_reports cr
		JOIN users u ON cr.reporter_id = u.id
		JOIN comments c ON cr.comment_id = c.id
		JOIN users cu ON c.author_id = cu.id
		ORDER BY cr.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourt et formate les résultats
	var reports []CommentReport
	for rows.Next() {
		var r CommentReport
		err := rows.Scan(&r.ID, &r.CommentID, &r.Reporter, &r.Reason, &r.CreatedAt, &r.CommentText, &r.CommentAuthor, &r.CommentAuthorID)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}

// Crée un nouveau signalement de commentaire
func CreateCommentReport(commentID int, reporterID, reason string) error {
	_, err := database.DB.Exec(`
		INSERT INTO comment_reports (comment_id, reporter_id, reason, created_at)
		VALUES (?, ?, ?, ?)`,
		commentID, reporterID, reason, time.Now(),
	)
	return err
}

// Supprime un signalement de commentaire par son ID
func DeleteCommentReportByID(reportID int) error {
	_, err := database.DB.Exec("DELETE FROM comment_reports WHERE id = ?", reportID)
	return err
}
