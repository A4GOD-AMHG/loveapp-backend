// Paquete services implementa la lógica de negocio de la aplicación.
package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
)

// MessageService define el contrato de operaciones de negocio sobre mensajes.
type MessageService interface {
	// SendMessage envía un nuevo mensaje del remitente al otro usuario del sistema.
	SendMessage(senderID uint, content string) (*models.Message, error)
	// EditMessage edita el contenido de un mensaje si fue enviado hace menos de 1 hora.
	EditMessage(userID uint, messageID int64, content string) (*models.Message, error)
	// DeleteMessage elimina un mensaje si fue enviado hace menos de 1 hora.
	DeleteMessage(userID uint, messageID int64) error
	// MarkAsRead marca un mensaje como leído por el destinatario.
	MarkAsRead(userID uint, messageID int64) error
	// MarkAsDelivered marca un mensaje como entregado al destinatario.
	MarkAsDelivered(userID uint, messageID int64) error
	// GetConversation retorna el historial paginado de mensajes entre dos usuarios.
	GetConversation(userID uint, page, perPage int) ([]models.Message, error)
	// GetUnreadCount retorna la cantidad de mensajes no leídos del usuario autenticado.
	GetUnreadCount(userID uint) (int, error)
}

// messageService es la implementación concreta de MessageService.
type messageService struct {
	messageRepo repository.MessageRepository         // Repositorio de mensajes
	userRepo    *repository.UserRepository           // Repositorio de usuarios (para encontrar al destinatario)
	deviceRepo  repository.DevicePushTokenRepository // Repositorio de tokens push del destinatario
	hub         *websocket.Hub                       // Hub de WebSocket para enviar eventos en tiempo real
	pushService PushService                          // Servicio para notificaciones push
}

// NewMessageService crea y retorna una nueva instancia de messageService
// con todas sus dependencias inyectadas.
func NewMessageService(
	messageRepo repository.MessageRepository,
	userRepo *repository.UserRepository,
	deviceRepo repository.DevicePushTokenRepository,
	pushService PushService,
	hub *websocket.Hub,
) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		deviceRepo:  deviceRepo,
		pushService: pushService,
		hub:         hub,
	}
}

// SendMessage crea y persiste un nuevo mensaje del remitente al otro usuario del sistema.
// Determina automáticamente el destinatario como el único usuario distinto al remitente.
// Difunde el evento "message_sent" por WebSocket tras guardar el mensaje.
func (s *messageService) SendMessage(senderID uint, content string) (*models.Message, error) {
	// Validar que el contenido del mensaje no esté vacío
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	// Obtener automáticamente al otro usuario como destinatario
	otherUser, err := s.userRepo.GetOtherUser(senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to determine receiver: %w", err)
	}

	// Construir la entidad mensaje con estado inicial "sent"
	msg := &models.Message{
		SenderID:   senderID,
		ReceiverID: uint(otherUser.ID),
		Content:    content,
		Status:     "sent",
	}

	// Persistir el mensaje en la base de datos
	id, err := s.messageRepo.Create(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	msg.ID = uint(id)

	// Recargar el mensaje completo desde la base de datos (con datos JOIN del remitente/destinatario)
	createdMsg, err := s.messageRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created message: %w", err)
	}

	// Notificar a los clientes conectados vía WebSocket
	s.hub.BroadcastEvent(&websocket.Event{
		Type:       websocket.EventMessageSent,
		SenderID:   createdMsg.SenderID,
		ReceiverID: createdMsg.ReceiverID,
		Payload:    *createdMsg,
	})

	receiverTokens, err := s.deviceRepo.FindByUserID(int64(createdMsg.ReceiverID))
	if err != nil {
		log.Printf("error fetching push tokens for user %d: %v", createdMsg.ReceiverID, err)
	} else if err := s.pushService.SendNewMessage(receiverTokens, models.PushMessagePayload{
		Type:       "chat_message",
		ChatID:     "private-main",
		MessageID:  createdMsg.ID,
		SenderID:   createdMsg.SenderID,
		SenderName: normalizeSenderName(createdMsg.Sender.Name),
		Content:    createdMsg.Content,
		CreatedAt:  createdMsg.CreatedAt,
	}); err != nil {
		log.Printf("error sending push notification for message %d: %v", createdMsg.ID, err)
	}

	return createdMsg, nil
}

func normalizeSenderName(rawName string) string {
	name := strings.TrimSpace(rawName)
	lowerName := strings.ToLower(name)

	switch {
	case strings.Contains(lowerName, "alexis"):
		return "Alexis"
	case strings.Contains(lowerName, "anyel"):
		return "Anyel"
	default:
		return name
	}
}

// EditMessage actualiza el contenido de un mensaje existente.
// Solo el remitente puede editar su propio mensaje y únicamente dentro de la primera hora de enviado.
// Difunde el evento "message_updated" por WebSocket tras la edición.
func (s *messageService) EditMessage(userID uint, messageID int64, content string) (*models.Message, error) {
	// Verificar que el mensaje existe
	msg, err := s.messageRepo.FindByID(messageID)
	if err != nil || msg == nil {
		return nil, errors.New("message not found")
	}

	// Verificar que el usuario solicitante es el remitente del mensaje
	if msg.SenderID != userID {
		return nil, errors.New("user not authorized to edit this message")
	}

	// Verificar que no ha pasado más de 1 hora desde el envío
	if time.Since(msg.CreatedAt) > time.Hour {
		return nil, errors.New("message can no longer be edited")
	}

	// Actualizar el contenido en la base de datos
	if err := s.messageRepo.UpdateContent(messageID, content); err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	// Recargar el mensaje actualizado para retornar el estado completo
	updatedMsg, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated message: %w", err)
	}

	// Notificar a los clientes conectados del cambio vía WebSocket
	s.hub.BroadcastEvent(&websocket.Event{
		Type:       websocket.EventMessageUpdated,
		SenderID:   updatedMsg.SenderID,
		ReceiverID: updatedMsg.ReceiverID,
		Payload: models.MessageStatusPayload{
			ID:        updatedMsg.ID,
			Status:    updatedMsg.Status,
			UpdatedAt: updatedMsg.UpdatedAt,
		},
	})

	return updatedMsg, nil
}

// DeleteMessage elimina permanentemente un mensaje de la base de datos.
// Solo el remitente puede eliminar su mensaje y únicamente dentro de la primera hora de enviado.
// Difunde el evento "message_deleted" por WebSocket tras la eliminación.
func (s *messageService) DeleteMessage(userID uint, messageID int64) error {
	// Verificar que el mensaje existe
	msg, err := s.messageRepo.FindByID(messageID)
	if err != nil || msg == nil {
		return errors.New("message not found")
	}

	// Verificar que el usuario solicitante es el remitente del mensaje
	if msg.SenderID != userID {
		return errors.New("user not authorized to delete this message")
	}

	// Verificar que no ha pasado más de 1 hora desde el envío
	if time.Since(msg.CreatedAt) > time.Hour {
		return errors.New("message can no longer be deleted")
	}

	// Eliminar el mensaje de la base de datos
	if err := s.messageRepo.Delete(messageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	// Notificar a los clientes conectados del evento de eliminación vía WebSocket
	// Se envían los datos del mensaje original para que el cliente pueda identificarlo
	s.hub.BroadcastEvent(&websocket.Event{
		Type:       websocket.EventMessageDeleted,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Payload:    *msg,
	})

	return nil
}

// MarkAsRead marca un mensaje como leído por el destinatario.
// Solo el destinatario puede marcar el mensaje como leído.
// Si el mensaje ya estaba marcado como leído, no hace nada.
// Difunde el evento "message_read" por WebSocket tras el cambio.
func (s *messageService) MarkAsRead(userID uint, messageID int64) error {
	// Verificar que el mensaje existe
	msg, err := s.messageRepo.FindByID(messageID)
	if err != nil || msg == nil {
		return errors.New("message not found")
	}

	// Solo el destinatario puede marcar el mensaje como leído
	if msg.ReceiverID != userID {
		return errors.New("user not authorized to mark this message as read")
	}

	// Evitar actualización innecesaria si ya está marcado como leído
	if msg.Status == "read" {
		return nil
	}

	// Actualizar el estado del mensaje en la base de datos
	if err := s.messageRepo.UpdateStatus(messageID, "read"); err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	// Actualizar el struct local para el evento WebSocket
	msg.Status = "read"
	msg.UpdatedAt = time.Now()

	// Notificar a los clientes conectados del cambio de estado vía WebSocket
	s.hub.BroadcastEvent(&websocket.Event{
		Type:       websocket.EventMessageRead,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Payload: models.MessageStatusPayload{
			ID:        msg.ID,
			Status:    msg.Status,
			UpdatedAt: msg.UpdatedAt,
		},
	})

	return nil
}

// MarkAsDelivered marca un mensaje como entregado al destinatario.
// Solo el destinatario puede marcar el mensaje como entregado.
// Si el mensaje ya está en estado "delivered" o "read", no hace nada.
// Difunde el evento "message_delivered" por WebSocket tras el cambio.
func (s *messageService) MarkAsDelivered(userID uint, messageID int64) error {
	// Verificar que el mensaje existe
	msg, err := s.messageRepo.FindByID(messageID)
	if err != nil || msg == nil {
		return errors.New("message not found")
	}

	// Solo el destinatario puede marcar el mensaje como entregado
	if msg.ReceiverID != userID {
		return errors.New("user not authorized to mark this message as delivered")
	}

	// Evitar regresión de estado: "read" no puede volver a "delivered"
	if msg.Status == "delivered" || msg.Status == "read" {
		return nil
	}

	// Actualizar el estado del mensaje en la base de datos
	if err := s.messageRepo.UpdateStatus(messageID, "delivered"); err != nil {
		return fmt.Errorf("failed to mark message as delivered: %w", err)
	}

	// Actualizar el struct local para el evento WebSocket
	msg.Status = "delivered"
	msg.UpdatedAt = time.Now()

	// Notificar a los clientes conectados del cambio de estado vía WebSocket
	s.hub.BroadcastEvent(&websocket.Event{
		Type:       websocket.EventMessageDelivered,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Payload: models.MessageStatusPayload{
			ID:        msg.ID,
			Status:    msg.Status,
			UpdatedAt: msg.UpdatedAt,
		},
	})

	return nil
}

// GetConversation retorna el historial de mensajes entre el usuario actual y el otro usuario,
// con soporte para paginación. Determina automáticamente al otro participante de la conversación.
func (s *messageService) GetConversation(userID uint, page, perPage int) ([]models.Message, error) {
	// Identificar al otro usuario de la conversación
	otherUser, err := s.userRepo.GetOtherUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to determine other user: %w", err)
	}
	// Delegar la consulta paginada al repositorio de mensajes
	return s.messageRepo.GetConversation(userID, uint(otherUser.ID), page, perPage)
}

// GetUnreadCount retorna la cantidad de mensajes no leídos cuyo destinatario es userID.
func (s *messageService) GetUnreadCount(userID uint) (int, error) {
	return s.messageRepo.CountUnreadByReceiver(userID)
}
