package handlers

import (
	"html/template"
	"net/http"

	"chat-platform/middleware"
	"chat-platform/models"
	"chat-platform/service"
)

type AuthHandler struct {
	AuthService    *service.AuthService
	SessionService *service.SessionService
	Templates      *template.Template
}

// Constructor
func NewAuthHandler(
	authService *service.AuthService,
	sessionService *service.SessionService,
	templates *template.Template,
) *AuthHandler {
	return &AuthHandler{
		AuthService:    authService,
		SessionService: sessionService,
		Templates:      templates,
	}
}

// Register displays the registration page and handles registration.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		if err := h.Templates.ExecuteTemplate(w, "register.html", nil); err != nil {
			http.Error(w, "failed to render page", http.StatusInternalServerError)
		}

	case http.MethodPost:

		user := &models.User{
			FirstName: r.FormValue("first_name"),
			LastName:  r.FormValue("last_name"),
			Username:  r.FormValue("username"),
			Email:     r.FormValue("email"),
		}

		password := r.FormValue("password")

		if err := h.AuthService.Register(user, password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Login displays the login page and authenticates users.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		if err := h.Templates.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, "failed to render page", http.StatusInternalServerError)
		}

	case http.MethodPost:

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := h.AuthService.Login(email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		session, err := h.SessionService.CreateSession(user.ID)
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.Token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400,
		})

		http.Redirect(w, r, "/chat", http.StatusSeeOther)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Logout signs the user out.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_token")
	if err == nil {
		_ = h.SessionService.Logout(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Profile displays the profile page and handles updates.
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {

	user, ok := middleware.GetUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {

	case http.MethodGet:

		data := struct {
			User interface{}
		}{
			User: user,
		}

		if err := h.Templates.ExecuteTemplate(w, "profile.html", data); err != nil {
			http.Error(w, "failed to render profile page", http.StatusInternalServerError)
		}

	case http.MethodPost:

		user.FirstName = r.FormValue("first_name")
		user.LastName = r.FormValue("last_name")
		user.Avatar = r.FormValue("avatar")

		if err := h.AuthService.UpdateUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/chat", http.StatusSeeOther)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}