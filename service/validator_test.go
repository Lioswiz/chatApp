package service

import (
	"testing"

	"chat-platform/models"
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name     string
		user     *models.User
		password string
		wantErr  bool
		errStr   string
	}{
		{
			name: "valid user",
			user: &models.User{
				FirstName: "John",
				LastName:  "Doe",
				Username:  "johndoe",
				Email:     "john@example.com",
			},
			password: "password123",
			wantErr:  false,
		},
		{
			name: "missing first name",
			user: &models.User{
				FirstName: "",
				LastName:  "Doe",
				Username:  "johndoe",
				Email:     "john@example.com",
			},
			password: "password123",
			wantErr:  true,
			errStr:   "first name is required",
		},
		{
			name: "missing last name",
			user: &models.User{
				FirstName: "John",
				LastName:  "  ",
				Username:  "johndoe",
				Email:     "john@example.com",
			},
			password: "password123",
			wantErr:  true,
			errStr:   "last name is required",
		},
		{
			name: "missing username",
			user: &models.User{
				FirstName: "John",
				LastName:  "Doe",
				Username:  "",
				Email:     "john@example.com",
			},
			password: "password123",
			wantErr:  true,
			errStr:   "username is required",
		},
		{
			name: "missing email",
			user: &models.User{
				FirstName: "John",
				LastName:  "Doe",
				Username:  "johndoe",
				Email:     "",
			},
			password: "password123",
			wantErr:  true,
			errStr:   "email is required",
		},
		{
			name: "short password",
			user: &models.User{
				FirstName: "John",
				LastName:  "Doe",
				Username:  "johndoe",
				Email:     "john@example.com",
			},
			password: "pass",
			wantErr:  true,
			errStr:   "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUser(tt.user, tt.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errStr {
				t.Errorf("ValidateUser() error string = %q, want %q", err.Error(), tt.errStr)
			}
		})
	}
}
