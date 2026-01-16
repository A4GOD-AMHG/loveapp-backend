package models

import "time"

// User representa un usuario en el sistema
// @Description Información de cuenta de usuario
type User struct {
	ID        int64     `json:"id" db:"id" example:"1"`
	Name      string    `json:"name" db:"name" example:"Anyel"`
	Username  string    `json:"username" db:"username" example:"anyel"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// LoginRequest representa la carga útil de solicitud de inicio de sesión
// @Description Credenciales de inicio de sesión
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"anyel"`    // Nombre de usuario
	Password string `json:"password" validate:"required" example:"password"` // Contraseña
}

// LoginResponse representa la respuesta de inicio de sesión
// @Description Respuesta de inicio de sesión con token e información de usuario
type LoginResponse struct {
	Message string `json:"message" example:"Inicio de sesión exitoso"`              // Mensaje de éxito
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Token JWT
	User    User   `json:"user"`                                                    // Información de usuario
}

// ChangePasswordRequest representa la solicitud de cambio de contraseña
// @Description Solicitud de cambio de contraseña
type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required" example:"newpassword"` // Nueva contraseña
}

// ChangePasswordResponse representa la respuesta de cambio de contraseña
// @Description Respuesta de cambio de contraseña
type ChangePasswordResponse struct {
	Message string `json:"message" example:"¡Contraseña actualizada con éxito! ✨"` // Mensaje de éxito
}
