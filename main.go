package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/GiaKhanh16/GoLangOfficial/config"
	"github.com/GiaKhanh16/GoLangOfficial/handlers"
)

func main() {
	// Load config and DB connection
	_, dbConn := config.LoadConfig()
	defer dbConn.Close(context.Background())

	fmt.Println("Connected to database!")

	// HTTP routes
	http.HandleFunc("/events", handlers.EventsHandler(dbConn))
	http.HandleFunc("/auth", handlers.UsersHandler(dbConn))
	http.HandleFunc("/messages", handlers.MessagesHandler(dbConn))
	http.HandleFunc("/ws", handlers.WebsocketHandler(dbConn))
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
