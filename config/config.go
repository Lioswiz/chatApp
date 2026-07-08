package config

// Config holds all the application's configuration settings.
type Config struct {
	// HTTP Server
	Port string

	// SQLite database location
	DatabasePath string

	// Folder where uploaded files are stored
	UploadDir string

	// Maximum upload size (in bytes)
	MaxUploadSize int64

	// Secret key used for session cookies
	SessionSecret string
}

// Load returns the application's configuration.
func Load() *Config {
	return &Config{
		// Server configuration
		Port: ":8080",

		// Database
		DatabasePath: "./database/chat.db",

		// Uploads
		UploadDir: "./static/uploads",

		// 50 MB
		MaxUploadSize: 50 * 1024 * 1024,

		// Change this before deploying to production
		SessionSecret: "my-super-secret-session-key",
	}
}
