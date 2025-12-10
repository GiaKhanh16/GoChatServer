package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GiaKhanh16/GoLangOfficial/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func UsersHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Insert if not exists
		_, err := conn.Exec(r.Context(),
			`INSERT INTO "user" (email) VALUES ($1) ON CONFLICT (email) DO NOTHING`,
			input.Email,
		)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			return
		}

		var userID uuid.UUID
		if err := conn.QueryRow(r.Context(),
			`SELECT id FROM "user" WHERE email = $1`,
			input.Email,
		).Scan(&userID); err != nil {
			http.Error(w, "Failed to fetch user ID", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.User{
			ID:    userID,
			Email: input.Email,
		})
	}
}
