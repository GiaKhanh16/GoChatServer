package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GiaKhanh16/GoLangOfficial/models"
	"github.com/jackc/pgx/v5"
)

func EventsHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := conn.Query(r.Context(), `SELECT id, event_name FROM event`)
		if err != nil {
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var events []models.Event
		for rows.Next() {
			var e models.Event
			if err := rows.Scan(&e.ID, &e.EventName); err != nil {
				http.Error(w, "Failed to scan event", http.StatusInternalServerError)
				return
			}
			events = append(events, e)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}
