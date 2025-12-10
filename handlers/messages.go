package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GiaKhanh16/GoLangOfficial/db"
	"github.com/jackc/pgx/v5"
)

func MessagesHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Fetch last messages
			eventID := r.URL.Query().Get("eventId")
			if eventID == "" {
				http.Error(w, "eventId is required", http.StatusBadRequest)
				return
			}

			messages, err := db.FetchLastFiveMessages(conn, eventID)
			if err != nil {
				http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(messages)

		case http.MethodPut:
			// Update a message
			id := r.URL.Query().Get("id")
			newText := r.URL.Query().Get("text")
			if id == "" || newText == "" {
				http.Error(w, "id and text are required", http.StatusBadRequest)
				return
			}

			if err := db.UpdateMessage(conn, id, newText); err != nil {
				http.Error(w, "Failed to update message", http.StatusInternalServerError)
				return
			}

			w.Write([]byte("message updated"))

		case http.MethodDelete:
			// Delete a message
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			if err := db.DeleteMessage(conn, id); err != nil {
				http.Error(w, "Failed to delete message", http.StatusInternalServerError)
				return
			}

			w.Write([]byte("message deleted"))

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
