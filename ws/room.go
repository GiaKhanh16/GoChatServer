package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/GiaKhanh16/GoLangOfficial/models"
	"github.com/coder/websocket"
)

// Room represents a chat room with multiple WebSocket clients
type Room struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan models.Message
	Mu        sync.Mutex
}

// NewRoom creates a new chat room
func NewRoom() *Room {
	r := &Room{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan models.Message, 100), // buffered channel
	}

	// Start the broadcast goroutine
	go r.run()
	return r
}

// run continuously broadcasts messages to all clients
func (r *Room) run() {
	for msg := range r.Broadcast {
		r.Mu.Lock()
		for client := range r.Clients {
			go func(c *websocket.Conn, m models.Message) {
				data, _ := json.Marshal(m)
				if err := c.Write(context.Background(), websocket.MessageText, data); err != nil {
					log.Println("Write error:", err)
					c.Close(websocket.StatusInternalError, "write error")

					r.Mu.Lock()
					delete(r.Clients, c)
					r.Mu.Unlock()
				}
			}(client, msg)
		}
		r.Mu.Unlock()
	}
}
