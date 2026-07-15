package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"chat-platform/models"
	"chat-platform/repository"
)

const (
	MaxUploadSize = 50 * 1024 * 1024 // 50 MB
)

var AllowedMimeTypes = map[string]bool{
	// Images
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,

	// Videos
	"video/mp4":  true,
	"video/webm": true,
	"video/ogg":  true,

	// Documents
	"application/pdf": true,
	"text/plain":      true,
}

type UploadService struct {
	UploadRepo *repository.UploadRepository
	UploadDir  string
}

// Constructor
func NewUploadService(uploadRepo *repository.UploadRepository, uploadDir string) *UploadService {
	return &UploadService{
		UploadRepo: uploadRepo,
		UploadDir:  uploadDir,
	}
}

// SaveFile validates, stores the file, and saves its metadata.
func (s *UploadService) SaveFile(file multipart.File, header *multipart.FileHeader, userID int) (*models.Upload, error) {

	defer file.Close()

	// Check file size
	if header.Size > MaxUploadSize {
		return nil, errors.New("file exceeds maximum size")
	}

	// Check MIME type
	contentType := header.Header.Get("Content-Type")
	if !AllowedMimeTypes[contentType] {
		return nil, errors.New("unsupported file type")
	}

	// Ensure upload directory exists
	err := os.MkdirAll(s.UploadDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Create a unique filename
	extension := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)

	path := filepath.Join(s.UploadDir, filename)

	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy file
	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	// Store relative path for URL consumption, replacing backslashes with forward slashes
	webPath := filepath.ToSlash(filepath.Join("/static/uploads", filename))

	upload := &models.Upload{
		UploadedBy: userID,
		FileName:   header.Filename,
		FilePath:   webPath,
		FileSize:   header.Size,
		MimeType:   strings.ToLower(contentType),
	}

	err = s.UploadRepo.SaveUpload(upload)
	if err != nil {
		return nil, err
	}

	return upload, nil
}