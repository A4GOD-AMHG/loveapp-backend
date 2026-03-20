// Paquete websocket implementa la comunicación en tiempo real mediante WebSockets.
package websocket

// Constantes que definen los tipos de eventos WebSocket emitidos por el servidor.
// El cliente debe escuchar estos tipos para actualizar la UI en tiempo real.
const (
	EventMessageSent      = "message_sent"      // Se emite cuando se envía un nuevo mensaje
	EventMessageUpdated   = "message_updated"   // Se emite cuando se edita el contenido de un mensaje
	EventMessageDeleted   = "message_deleted"   // Se emite cuando se elimina un mensaje
	EventMessageRead      = "message_read"      // Se emite cuando el destinatario marca el mensaje como leído
	EventMessageDelivered = "message_delivered" // Se emite cuando el mensaje es entregado al destinatario
)

// Event representa un evento WebSocket enviado desde el servidor a los clientes conectados.
// Contiene el tipo de evento, los participantes relevantes y el payload serializable.
type Event struct {
	Type       string `json:"type"`    // Tipo del evento (ver constantes EventMessage*)
	Payload    any    `json:"payload"` // Datos serializables del evento
	SenderID   uint   `json:"-"`       // Destinatario lógico para el emisor del evento
	ReceiverID uint   `json:"-"`       // Destinatario lógico para el receptor del evento
}
