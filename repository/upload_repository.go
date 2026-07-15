package repository

import (
	"database/sql"

	"chat-platform/models"
)

type UploadRepository struct {
	DB *sql.DB
}

func NewUploadRepository(db *sql.DB) *UploadRepository {
	return &UploadRepository{
		DB: db,
	}
}

// SaveUpload stores metadata about an uploaded file.
func (r *UploadRepository) SaveUpload(upload *models.Upload) error {
	query := `
	INSERT INTO uploads
	(uploaded_by, file_name, file_path, file_size, mime_type)
	VALUES (?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(
		query,
		upload.UploadedBy,
		upload.FileName,
		upload.FilePath,
		upload.FileSize,
		upload.MimeType,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	upload.ID = int(id)
	return nil
}

// GetUpload returns an upload by ID.
func (r *UploadRepository) GetUpload(id int) (*models.Upload, error) {
	query := `
	SELECT
		id,
		uploaded_by,
		file_name,
		file_path,
		file_size,
		mime_type
	FROM uploads
	WHERE id = ?`

	upload := &models.Upload{}

	err := r.DB.QueryRow(query, id).Scan(
		&upload.ID,
		&upload.UploadedBy,
		&upload.FileName,
		&upload.FilePath,
		&upload.FileSize,
		&upload.MimeType,
	)

	if err != nil {
		return nil, err
	}

	return upload, nil
}