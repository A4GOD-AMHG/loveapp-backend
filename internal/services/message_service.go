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

const (
	StatusSent      = "sent"
	StatusDelivered = "delivered"
	StatusRead      = "read"

	// editWindow es el tiempo durante el cual el remitente puede editar o borrar.
	editWindow = time.Hour
)

// statusRank ordena los estados para impedir regresiones (read -> delivered).
var statusRank = map[string]int{StatusSent: 0, StatusDelivered: 1, StatusRead: 2}

// MessageService define el contrato de operaciones de negocio sobre mensajes.
type MessageService interface {
	SendMessage(senderID uint, content string) (*models.Message, error)
	EditMessage(userID uint, messageID int64, content string) (*models.Message, error)
	DeleteMessage(userID uint, messageID int64) error
	MarkAsRead(userID uint, messageID int64) error
	MarkAsDelivered(userID uint, messageID int64) error
	GetConversation(userID uint, page, perPage int) ([]models.Message, error)
	GetUnreadCount(userID uint) (int, error)
}

// EventBroadcaster difunde eventos en tiempo real (implementado por websocket.Hub).
type EventBroadcaster interface {
	BroadcastEvent(event *websocket.Event)
}

// OtherUserFinder resuelve al otro integrante de la pareja.
type OtherUserFinder interface {
	GetOtherUser(currentUserID uint) (*models.User, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
	userRepo    OtherUserFinder
	deviceRepo  repository.DevicePushTokenRepository
	hub         EventBroadcaster
	pushService PushService
}

// NewMessageService crea un servicio de mensajes con sus dependencias inyectadas.
func NewMessageService(
	messageRepo repository.MessageRepository,
	userRepo OtherUserFinder,
	deviceRepo repository.DevicePushTokenRepository,
	pushService PushService,
	hub EventBroadcaster,
) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		deviceRepo:  deviceRepo,
		pushService: pushService,
		hub:         hub,
	}
}

// SendMessage persiste un mensaje hacia el otro usuario, lo difunde por WebSocket
// y envía la notificación push en segundo plano.
func (s *messageService) SendMessage(senderID uint, content string) (*models.Message, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ErrInvalidInput
	}

	otherUser, err := s.otherUser(senderID)
	if err != nil {
		return nil, err
	}

	id, err := s.messageRepo.Create(&models.Message{
		SenderID:   senderID,
		ReceiverID: uint(otherUser.ID),
		Content:    content,
		Status:     StatusSent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	msg, err := s.messageRepo.FindByID(id)
	if err != nil || msg == nil {
		return nil, fmt.Errorf("failed to fetch created message: %v", err)
	}

	s.broadcast(websocket.EventMessageSent, msg, *msg)
	// El push no debe retrasar la respuesta HTTP (FCM puede tardar segundos).
	go s.notifyReceiver(*msg)

	return msg, nil
}

func (s *messageService) notifyReceiver(msg models.Message) {
	tokens, err := s.deviceRepo.FindByUserID(int64(msg.ReceiverID))
	if err != nil {
		log.Printf("error fetching push tokens for user %d: %v", msg.ReceiverID, err)
		return
	}
	if err := s.pushService.SendNewMessage(tokens, models.PushMessagePayload{
		Type:       "chat_message",
		ChatID:     "private-main",
		MessageID:  msg.ID,
		SenderID:   msg.SenderID,
		SenderName: normalizeSenderName(msg.Sender.Name),
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt,
	}); err != nil {
		log.Printf("error sending push notification for message %d: %v", msg.ID, err)
	}
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

// EditMessage actualiza el contenido; solo el remitente y dentro de editWindow.
func (s *messageService) EditMessage(userID uint, messageID int64, content string) (*models.Message, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ErrInvalidInput
	}

	if _, err := s.findOwnRecent(userID, messageID); err != nil {
		return nil, err
	}

	if err := s.messageRepo.UpdateContent(messageID, content); err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	updated, err := s.messageRepo.FindByID(messageID)
	if err != nil || updated == nil {
		return nil, fmt.Errorf("failed to fetch updated message: %v", err)
	}

	// Se envía el mensaje completo para que el cliente pueda refrescar el contenido.
	s.broadcast(websocket.EventMessageUpdated, updated, *updated)
	return updated, nil
}

// DeleteMessage elimina un mensaje; solo el remitente y dentro de editWindow.
func (s *messageService) DeleteMessage(userID uint, messageID int64) error {
	msg, err := s.findOwnRecent(userID, messageID)
	if err != nil {
		return err
	}

	if err := s.messageRepo.Delete(messageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	s.broadcast(websocket.EventMessageDeleted, msg, *msg)
	return nil
}

// MarkAsRead marca un mensaje como leído; solo el destinatario puede hacerlo.
func (s *messageService) MarkAsRead(userID uint, messageID int64) error {
	return s.advanceStatus(userID, messageID, StatusRead, websocket.EventMessageRead)
}

// MarkAsDelivered marca un mensaje como entregado; solo el destinatario puede hacerlo.
func (s *messageService) MarkAsDelivered(userID uint, messageID int64) error {
	return s.advanceStatus(userID, messageID, StatusDelivered, websocket.EventMessageDelivered)
}

func (s *messageService) advanceStatus(userID uint, messageID int64, status, eventType string) error {
	msg, err := s.find(messageID)
	if err != nil {
		return err
	}
	if msg.ReceiverID != userID {
		return ErrForbidden
	}
	if statusRank[msg.Status] >= statusRank[status] {
		return nil
	}

	if err := s.messageRepo.UpdateStatus(messageID, status); err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	msg.Status = status
	msg.UpdatedAt = time.Now()
	s.broadcast(eventType, msg, models.MessageStatusPayload{
		ID:        msg.ID,
		Status:    msg.Status,
		UpdatedAt: msg.UpdatedAt,
	})
	return nil
}

// GetConversation retorna el historial paginado con el otro usuario.
func (s *messageService) GetConversation(userID uint, page, perPage int) ([]models.Message, error) {
	otherUser, err := s.otherUser(userID)
	if err != nil {
		return nil, err
	}
	return s.messageRepo.GetConversation(userID, uint(otherUser.ID), page, perPage)
}

// GetUnreadCount retorna la cantidad de mensajes no leídos cuyo destinatario es userID.
func (s *messageService) GetUnreadCount(userID uint) (int, error) {
	return s.messageRepo.CountUnreadByReceiver(userID)
}

func (s *messageService) otherUser(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetOtherUser(userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNoReceiver
	}
	if err != nil {
		return nil, fmt.Errorf("failed to determine receiver: %w", err)
	}
	return user, nil
}

func (s *messageService) find(messageID int64) (*models.Message, error) {
	msg, err := s.messageRepo.FindByID(messageID)
	if err != nil || msg == nil {
		return nil, ErrNotFound
	}
	return msg, nil
}

// findOwnRecent valida que el mensaje exista, sea del usuario y siga editable.
func (s *messageService) findOwnRecent(userID uint, messageID int64) (*models.Message, error) {
	msg, err := s.find(messageID)
	if err != nil {
		return nil, err
	}
	if msg.SenderID != userID {
		return nil, ErrForbidden
	}
	if time.Since(msg.CreatedAt) > editWindow {
		return nil, ErrExpired
	}
	return msg, nil
}

func (s *messageService) broadcast(eventType string, msg *models.Message, payload any) {
	s.hub.BroadcastEvent(&websocket.Event{
		Type:       eventType,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Payload:    payload,
	})
}
