package models

import (
	"Forum/internal/database"
	"time"
)

type CommentReport struct {
	ID                  int
	CommentID           int
	Reporter            string
	Reason              string
	CreatedAt           time.Time
	CommentText         string
	CommentAuthor       string
	CommentAuthorID     string
}

func GetAllCommentReports() ([]CommentReport, error) {
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

func CreateCommentReport(commentID int, reporterID, reason string) error {
	_, err := database.DB.Exec(`
		INSERT INTO comment_reports (comment_id, reporter_id, reason)
		VALUES (?, ?, ?)`,
		commentID, reporterID, reason,
	)
	return err
}

func DeleteCommentReportByID(reportID int) error {
	_, err := database.DB.Exec("DELETE FROM comment_reports WHERE id = ?", reportID)
	return err
}
