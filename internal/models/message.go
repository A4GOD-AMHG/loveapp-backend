// Paquete models define las estructuras de datos utilizadas en toda la aplicación.
package models

import (
	"time"
)

// Message representa un mensaje enviado entre los dos usuarios del sistema.
type Message struct {
	ID         uint        `json:"id"`          // Identificador único del mensaje
	SenderID   uint        `json:"sender_id"`   // ID del usuario que envió el mensaje
	Sender     UserSummary `json:"sender"`      // Información resumida del remitente
	ReceiverID uint        `json:"receiver_id"` // ID del usuario que recibe el mensaje
	Receiver   UserSummary `json:"receiver"`    // Información resumida del destinatario
	Content    string      `json:"content"`     // Contenido textual del mensaje
	Status     string      `json:"status"`      // Estado del mensaje: "sent", "delivered" o "read"
	CreatedAt  time.Time   `json:"created_at"`  // Fecha y hora de creación del mensaje
	UpdatedAt  time.Time   `json:"updated_at"`  // Fecha y hora de la última actualización
}

// MessageStatusPayload es el payload compacto para eventos de estado del mensaje.
type MessageStatusPayload struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PushMessagePayload contiene el cuerpo útil enviado al proveedor push.
type PushMessagePayload struct {
	Type       string    `json:"type"`
	MessageID  uint      `json:"message_id"`
	SenderID   uint      `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
