// Paquete models define las estructuras de datos utilizadas en toda la aplicación.
package models

import "time"

// User representa un usuario registrado en el sistema.
// La contraseña nunca se serializa en las respuestas JSON (json:"-").
// @Description Información de cuenta de usuario
type User struct {
	ID        int64     `json:"id" db:"id" example:"1"`                                    // Identificador único del usuario
	Name      string    `json:"name" db:"name" example:"Anyel"`                            // Nombre real del usuario
	Username  string    `json:"username" db:"username" example:"anyel"`                    // Nombre de usuario único para autenticación
	Password  string    `json:"-" db:"password"`                                           // Hash de la contraseña (nunca expuesto en respuestas)
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2024-01-01T00:00:00Z"` // Fecha de creación de la cuenta
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2024-01-01T00:00:00Z"` // Fecha de la última actualización
}

// UserSummary representa una versión reducida del usuario para respuestas anidadas o públicas.
// Se usa dentro de otras estructuras como Message para mostrar info del remitente/destinatario.
// @Description Información resumida de usuario
type UserSummary struct {
	ID       int64  `json:"id" example:"1"`           // Identificador único del usuario
	Name     string `json:"name" example:"Anyel"`     // Nombre real del usuario
	Username string `json:"username" example:"anyel"` // Nombre de usuario
}

// LoginRequest contiene las credenciales enviadas para autenticarse en el sistema.
// @Description Credenciales de inicio de sesión
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"anyel"` // Nombre de usuario (obligatorio)
	Password string `json:"password" validate:"required" example:"anyel"` // Contraseña en texto plano (obligatorio)
}

// LoginResponse contiene el token JWT y los datos del usuario tras un inicio de sesión exitoso.
// @Description Respuesta de inicio de sesión con token e información de usuario
type LoginResponse struct {
	Message     string `json:"message" example:"Inicio de sesión exitoso"`              // Mensaje de confirmación
	Token       string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Token JWT para autenticación de solicitudes
	User        User   `json:"user"`                                                    // Datos del usuario autenticado (sin contraseña)
	UnreadCount int    `json:"unread_count" example:"3"`                                // Cantidad de mensajes no leídos para pintar el badge inicial
}

// ChangePasswordRequest contiene la nueva contraseña para el cambio de contraseña.
// @Description Solicitud de cambio de contraseña
type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required" example:"newpassword"` // Nueva contraseña en texto plano (mínimo 6 caracteres)
}

// ChangePasswordResponse es la respuesta retornada tras cambiar la contraseña exitosamente.
// @Description Respuesta de cambio de contraseña
type ChangePasswordResponse struct {
	Message string `json:"message" example:"¡Contraseña actualizada con éxito! ✨"` // Mensaje de confirmación del cambio
}

// RegisterPushTokenRequest contiene los datos necesarios para registrar un dispositivo push.
type RegisterPushTokenRequest struct {
	Platform   string `json:"platform" example:"android"`
	PushToken  string `json:"push_token" example:"FCM_OR_APNS_TOKEN"`
	DeviceName string `json:"device_name" example:"Pixel 8"`
}

// DeletePushTokenRequest contiene el token a desvincular del usuario autenticado.
type DeletePushTokenRequest struct {
	PushToken string `json:"push_token" example:"FCM_OR_APNS_TOKEN"`
}

// DevicePushToken representa un token push persistido por usuario y dispositivo.
type DevicePushToken struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Platform   string    `json:"platform"`
	PushToken  string    `json:"push_token"`
	DeviceName string    `json:"device_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UnreadCountResponse representa la cantidad total de mensajes no leídos.
type UnreadCountResponse struct {
	UnreadCount int `json:"unread_count" example:"7"`
}
