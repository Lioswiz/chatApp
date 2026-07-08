package routes

import (
	"database/sql"
	"html/template"
	"net/http"

	"chat-platform/handlers"
	"chat-platform/middleware"
	"chat-platform/repository"
	"chat-platform/service"
	"chat-platform/websocket"
)

// Register creates and registers all application routes.
func Register(db *sql.DB, hub *websocket.Hub) *http.ServeMux {

	// =====================================
	// Templates
	// =====================================

	templates := template.Must(template.ParseGlob("templates/*.html"))

	// =====================================
	// Repositories
	// =====================================

	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// =====================================
	// Services
	// =====================================

	authService := service.NewAuthService(userRepo)
	chatService := service.NewChatService(messageRepo)
	sessionService := service.NewSessionService(
		sessionRepo,
		userRepo,
	)

	// =====================================
	// HTTP Handlers
	// =====================================

	authHandler := handlers.NewAuthHandler(
		authService,
		sessionService,
		templates,
	)

	chatHandler := handlers.NewChatHandler(
		templates,
		chatService,
	)

	// =====================================
	// WebSocket Handler
	// =====================================

	wsHandler := websocket.NewWebSocketHandler(
		hub,
		chatService,
	)

	// =====================================
	// Router
	// =====================================

	mux := http.NewServeMux()

	// Home
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	// Authentication
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/logout", authHandler.Logout)

	// Chat Page (Protected)
	mux.Handle(
		"/chat",
		middleware.SessionMiddleware(
			sessionService,
			http.HandlerFunc(chatHandler.ChatPage),
		),
	)

	// WebSocket Endpoint (Protected)
	mux.Handle(
		"/ws",
		middleware.SessionMiddleware(
			sessionService,
			http.HandlerFunc(wsHandler.ServeWS),
		),
	)

	return mux
}