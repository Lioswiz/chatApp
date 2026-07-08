package handlers

import (
	"net/http"
)

type UploadHandler struct{}

// Upload handles file uploads.
// We'll implement the actual upload logic later.
func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.Error(w, "upload handler not implemented yet", http.StatusNotImplemented)
}