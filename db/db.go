package db

import (
	"context"
	"fmt"
	"log"

	"github.com/GiaKhanh16/GoLangOfficial/models"
	"github.com/jackc/pgx/v5"
)

// SaveMessage inserts a message into the database
func SaveMessage(conn *pgx.Conn, msg models.Message) error {
	_, err := conn.Exec(context.Background(),
		`INSERT INTO messages (id, user_id, user_name, event_id, content, email, created_at, image_name)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		msg.ID, msg.UserID, msg.UserName, msg.EventID, msg.Text, msg.Email, msg.CreatedAt, msg.ImageName,
	)
	if err != nil {
		log.Println("Error saving message:", err)
	}
	return err
}

// SaveReaction inserts a reaction into the database
func SaveReaction(conn *pgx.Conn, r models.Reaction) error {
	_, err := conn.Exec(context.Background(),
		`INSERT INTO message_reactions (id, message_id, user_id, emoji, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		r.ID, r.MessageID, r.UserID, r.Emoji, r.CreatedAt,
	)
	if err != nil {
		log.Println("Error saving reaction:", err)
	}
	return err
}

// FetchLastFiveMessages fetches the last 5 messages for a given event, oldest first
func FetchLastFiveMessages(conn *pgx.Conn, eventID string) ([]models.Message, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, user_id, user_name, event_id, content, email, created_at, image_name
		 FROM messages
		 WHERE event_id = $1
		 ORDER BY created_at DESC
		 LIMIT 5`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.UserName,
			&m.EventID,
			&m.Text,
			&m.Email,
			&m.CreatedAt,
			&m.ImageName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, m)
	}

	// Reverse slice so oldest messages come first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func UpdateMessage(conn *pgx.Conn, id string, newText string) error {
	_, err := conn.Exec(context.Background(), `UPDATE messages SET content = $1 WHERE id = $2`, newText, id)
	if err != nil {
		log.Println("Error updating message:", err)
	}
	return err
}

func DeleteMessage(conn *pgx.Conn, id string) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM messages
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		log.Println("Error deleting message:", err)
	}
	return err
}
