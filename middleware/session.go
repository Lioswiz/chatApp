package middleware

import (
	"context"
	"net/http"

	"chat-platform/models"
	"chat-platform/service"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

// SessionMiddleware authenticates requests using the session cookie.
func SessionMiddleware(sessionService *service.SessionService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := sessionService.ValidateSession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser returns the authenticated user stored in the request context.
func GetUser(r *http.Request) (*models.User, bool) {

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		return nil, false
	}

	return user, true
}

// GetUserID returns the authenticated user's ID.
func GetUserID(r *http.Request) (int, bool) {

	user, ok := GetUser(r)
	if !ok {
		return 0, false
	}

	return user.ID, true
}