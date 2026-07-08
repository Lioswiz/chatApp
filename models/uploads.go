package models

// Upload represents a file attached to a message.
type Upload struct {
	ID         int
	MessageID  int
	FileName   string
	FilePath   string
	FileSize   int64
	MimeType   string
}