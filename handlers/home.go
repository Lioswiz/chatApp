package handlers

import "net/http"

type HomeHandler struct{}

// Home redirects users to the login page.
func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}