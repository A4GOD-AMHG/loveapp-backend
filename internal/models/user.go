package models

import "time"

// User represents a user in the system
// @Description User account information
type User struct {
	ID        int64     `json:"id" db:"id" example:"1"`                                    // User ID
	Username  string    `json:"username" db:"username" example:"anyel"`                    // Username
	Password  string    `json:"-" db:"password"`                                           // Password (hidden in JSON)
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// LoginRequest represents the login request payload
// @Description Login credentials
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"anyel"`    // Username
	Password string `json:"password" validate:"required" example:"password"` // Password
}

// LoginResponse represents the login response
// @Description Login response with token and user info
type LoginResponse struct {
	Message string `json:"message" example:"Inicio de sesión exitoso"`              // Success message
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT token
	User    User   `json:"user"`                                                    // User information
}

// ChangePasswordRequest represents the change password request payload
// @Description Change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"oldpassword"` // Current password
	NewPassword string `json:"new_password" validate:"required" example:"newpassword"` // New password
}

// ChangePasswordResponse represents the change password response
// @Description Change password response
type ChangePasswordResponse struct {
	Message string `json:"message" example:"Contraseña cambiada exitosamente"` // Success message
}
