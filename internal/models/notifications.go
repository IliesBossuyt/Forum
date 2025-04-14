package models

import (
	"Forum/internal/database"
	"time"
)

type Notification struct {
	ID             int
	RecipientID    string
	SenderID       string
	SenderUsername string
	Type           string
	PostID         *int
	CommentID      *int
	Message        string
	Seen           bool
	CreatedAt      time.Time
}

func CreateNotification(n Notification) error {
	_, err := database.DB.Exec(`
        INSERT INTO notifications (recipient_id, sender_id, type, post_id, comment_id, message)
        VALUES (?, ?, ?, ?, ?, ?)`,
		n.RecipientID, n.SenderID, n.Type, n.PostID, n.CommentID, n.Message)
	return err
}

func GetNotificationsByUser(userID string) ([]Notification, error) {
	rows, err := database.DB.Query(`
        SELECT n.id, n.recipient_id, n.sender_id, u.username, n.type,
               n.post_id, n.comment_id, n.seen, n.created_at
        FROM notifications n
        JOIN users u ON n.sender_id = u.id
        WHERE n.recipient_id = ?
        ORDER BY n.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []Notification
	for rows.Next() {
		var n Notification
		var postID, commentID *int

		err := rows.Scan(&n.ID, &n.RecipientID, &n.SenderID, &n.SenderUsername, &n.Type,
			&postID, &commentID, &n.Seen, &n.CreatedAt)
		if err != nil {
			return nil, err
		}

		n.PostID = postID
		n.CommentID = commentID

		// Construire dynamiquement le message
		switch n.Type {
		case "like":
			n.Message = n.SenderUsername + " a aimé votre post."
		case "like_comment":
			n.Message = n.SenderUsername + " a aimé votre commentaire."
		case "comment":
			n.Message = n.SenderUsername + " a commenté votre post."
		case "warn":
			n.Message = n.SenderUsername + " vous a averti."
		default:
			n.Message = n.SenderUsername + " vous a envoyé une notification."
		}

		notifs = append(notifs, n)
	}

	return notifs, nil
}

func MarkNotificationsAsSeen(userID string) error {
	_, err := database.DB.Exec("UPDATE notifications SET seen = TRUE WHERE recipient_id = ?", userID)
	return err
}
