package main

import (
	"log"
	"net/http"

	"chat-platform/config"
	"chat-platform/database"
	"chat-platform/routes"
	"chat-platform/websocket"
)

func main() {
	// ==========================
	// Load Application Config
	// ==========================
	cfg := config.Load()

	// ==========================
	// Connect Database
	// ==========================
	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatal("Database Error:", err)
	}
	defer db.Close()

	log.Println("✅ Database Connected")

	// ==========================
	// Create WebSocket Hub
	// ==========================
	hub := websocket.NewHub()

	// Run the hub in the background
	go hub.Run()

	log.Println("✅ WebSocket Hub Running")

	// ==========================
	// Register All Routes
	// ==========================
	handler := routes.Register(db, hub)

	// ==========================
	// Start Server
	// ==========================
	log.Printf("🚀 Server running at http://localhost%s\n", cfg.Port)

	err = http.ListenAndServe(cfg.Port, handler)
	if err != nil {
		log.Fatal("Server Error:", err)
	}
}
