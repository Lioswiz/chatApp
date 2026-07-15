package models

// Upload represents a file attached to a message.
type Upload struct {
	ID         int
	UploadedBy int
	FileName   string
	FilePath   string
	FileSize   int64
	MimeType   string
}