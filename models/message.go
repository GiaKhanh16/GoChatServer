package models

import "time"

type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	EventID   string    `json:"event_id"`
	Email     string    `json:"email"`
	Text      string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ImageName string    `json:"image_name"`
}
