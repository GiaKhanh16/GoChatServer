package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/GiaKhanh16/GoLangOfficial/db"
	"github.com/GiaKhanh16/GoLangOfficial/models"
	"github.com/GiaKhanh16/GoLangOfficial/ws"
	"github.com/coder/websocket"
	"github.com/jackc/pgx/v5"
)

var rooms = make(map[string]*ws.Room)
var roomsMu sync.Mutex

func getRoom(eventId string) *ws.Room {
	roomsMu.Lock()
	defer roomsMu.Unlock()
	if rooms[eventId] == nil {
		rooms[eventId] = ws.NewRoom()
	}
	return rooms[eventId]
}

func WebsocketHandler(connDB *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventId := r.URL.Query().Get("eventId")
		if eventId == "" {
			http.Error(w, "eventId is required", http.StatusBadRequest)
			return
		}

		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"*"},
		})
		if err != nil {
			log.Println("Accept error:", err)
			return
		}

		log.Printf("Client connected to room %s", eventId)
		room := getRoom(eventId)

		room.Mu.Lock()
		room.Clients[conn] = true
		room.Mu.Unlock()

		defer func() {
			log.Printf("Client left room %s", eventId)
			room.Mu.Lock()
			delete(room.Clients, conn)
			room.Mu.Unlock()
			conn.Close(websocket.StatusNormalClosure, "bye")
		}()

		ctx := context.Background()
		for {
			msgType, data, err := conn.Read(ctx)
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			if msgType != websocket.MessageText {
				continue
			}
			log.Println("Received raw JSON:", string(data))

			var raw map[string]interface{}
			if err := json.Unmarshal(data, &raw); err != nil {
				log.Println("JSON unmarshal error:", err)
				continue
			}

			switch raw["type"] {
			case "reaction":
				var reaction models.Reaction
				if err := json.Unmarshal(data, &reaction); err != nil {
					log.Println("Reaction unmarshal error:", err)
					continue
				}

				if err := db.SaveReaction(connDB, reaction); err != nil {
					log.Println("Failed to save reaction:", err)
				}

				// broadcast reaction
				room.Mu.Lock()
				for client := range room.Clients {
					go client.Write(ctx, websocket.MessageText, data)
				}
				room.Mu.Unlock()

			default:
				var msg models.Message
				if err := json.Unmarshal(data, &msg); err != nil {
					log.Println("Message unmarshal error:", err)
					continue
				}

				if err := db.SaveMessage(connDB, msg); err != nil {
					log.Println("Failed to save message:", err)
				}

				room.Broadcast <- msg
			}
		}
	}
}
