package models

import "github.com/google/uuid"

type Event struct {
	ID        uuid.UUID `json:"id"`
	EventName string    `json:"event_name"`
}
