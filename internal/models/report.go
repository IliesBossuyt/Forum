package models

import (
	"Forum/internal/database"
	"time"
)

type Report struct {
	ID               int
	PostID           int
	Reporter         string
	Reason           string
	CreatedAt        time.Time
	PostContent      string
	PostAuthor       string
	PostAuthorID     string
	PostAuthorBanned bool
	PostImage        *string // Pointeur pour accepter NULL
}

func CreateReport(postID int, reporterID string, reason string) error {
	_, err := database.DB.Exec(`INSERT INTO reports (post_id, reporter_id, reason) VALUES (?, ?, ?)`,
		postID, reporterID, reason,
	)
	return err
}

func GetAllReports() ([]Report, error) {
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

func DeleteReportByID(reportID int) error {
	_, err := database.DB.Exec("DELETE FROM reports WHERE id = ?", reportID)
	return err
}
