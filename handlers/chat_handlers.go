package handlers

import (
	"html/template"
	"net/http"

	"chat-platform/middleware"
	"chat-platform/service"
)

type ChatHandler struct {
	Templates  *template.Template
	ChatService *service.ChatService
}

func NewChatHandler(
	templates *template.Template,
	chatService *service.ChatService,
) *ChatHandler {
	return &ChatHandler{
		Templates:   templates,
		ChatService: chatService,
	}
}

// ChatPage renders the main chat page.
func (h *ChatHandler) ChatPage(w http.ResponseWriter, r *http.Request) {

	// Get authenticated user from the session middleware.
	user, ok := middleware.GetUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Load all public chat messages.
	messages, err := h.ChatService.GetPublicMessages()
	if err != nil {
		http.Error(w, "failed to load messages", http.StatusInternalServerError)
		return
	}

	data := struct {
		User     interface{}
		Messages interface{}
	}{
		User:     user,
		Messages: messages,
	}

	err = h.Templates.ExecuteTemplate(w, "chat.html", data)
	if err != nil {
		http.Error(w, "failed to render chat page", http.StatusInternalServerError)
		return
	}
}