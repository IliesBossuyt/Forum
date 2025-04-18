package models

import (
	"Forum/internal/database"
	"time"
)

// Structure d'une notification
type Notification struct {
	ID             int       // Identifiant unique
	RecipientID    string    // ID du destinataire
	SenderID       string    // ID de l'émetteur
	SenderUsername string    // Nom de l'émetteur
	Type           string    // Type de notification
	PostID         *int      // ID du post concerné
	CommentID      *int      // ID du commentaire concerné
	Message        string    // Message de la notification
	Seen           bool      // État de lecture
	CreatedAt      time.Time // Date de création
}

// Crée une nouvelle notification
func CreateNotification(n Notification) error {
	_, err := database.DB.Exec(`
        INSERT INTO notifications (recipient_id, sender_id, type, post_id, comment_id, message, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		n.RecipientID, n.SenderID, n.Type, n.PostID, n.CommentID, n.Message, time.Now())
	return err
}

// Récupère les notifications d'un utilisateur
func GetNotificationsByUser(userID string) ([]Notification, error) {
	// Requête pour obtenir les notifications avec les infos de l'émetteur
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

	// Parcourt et formate les résultats
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

		// Construit le message selon le type
		switch n.Type {
		case "like_post":
			n.Message = n.SenderUsername + " a liké votre post."
		case "dislike_post":
			n.Message = n.SenderUsername + " a disliké votre post."
		case "like_comment":
			n.Message = n.SenderUsername + " a liké votre commentaire."
		case "dislike_comment":
			n.Message = n.SenderUsername + " a disliké votre commentaire."
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

// Marque toutes les notifications comme lues
func MarkNotificationsAsSeen(userID string) error {
	_, err := database.DB.Exec("UPDATE notifications SET seen = TRUE WHERE recipient_id = ?", userID)
	return err
}

// Supprime toutes les notifications d'un utilisateur
func DeleteAllNotificationsForUser(userID string) error {
	_, err := database.DB.Exec("DELETE FROM notifications WHERE recipient_id = ?", userID)
	return err
}
