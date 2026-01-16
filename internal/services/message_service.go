package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
)

type MessageService interface {
	SendMessage(senderID, receiverID uint, content string) (*models.Message, error)
	EditMessage(userID uint, messageID int64, content string) (*models.Message, error)
	DeleteMessage(userID uint, messageID int64) error
	GetConversation(user1ID, user2ID uint) ([]models.Message, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
	userRepo    *repository.UserRepository
	hub         *websocket.Hub
}

func NewMessageService(messageRepo repository.MessageRepository, userRepo *repository.UserRepository, hub *websocket.Hub) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		hub:         hub,
	}
}

func (s *messageService) SendMessage(senderID, receiverID uint, content string) (*models.Message, error) {
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	msg := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Status:     "sent",
	}

	id, err := s.messageRepo.Create(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	msg.ID = uint(id)
	// Refetch to get all fields populated by DB
	createdMsg, err := s.messageRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created message: %w", err)
	}

	// Broadcast event via WebSocket
	s.hub.BroadcastEvent(&websocket.Event{
		Type:    websocket.EventMessageSent,
		Payload: *createdMsg,
	})

	return createdMsg, nil
}

func (s *messageService) EditMessage(userID uint, messageID int64, content string) (*models.Message, error) {
	msg, err := s.messageRepo.FindByID(messageID)
	if err != nil || msg == nil {
		return nil, errors.New("message not found")
	}

	if msg.SenderID != userID {
		return nil, errors.New("user not authorized to edit this message")
	}

	if time.Since(msg.CreatedAt) > time.Hour {
		return nil, errors.New("message can no longer be edited")
	}

	if err := s.messageRepo.UpdateContent(messageID, content); err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	updatedMsg, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated message: %w", err)
	}

	// Broadcast event via WebSocket
	s.hub.BroadcastEvent(&websocket.Event{
		Type:    websocket.EventMessageUpdated,
		Payload: *updatedMsg,
	})

	return updatedMsg, nil
}

func (s *messageService) DeleteMessage(userID uint, messageID int64) error {
	msg, err := s.messageRepo.FindByID(messageID)
	if err != nil || msg == nil {
		return errors.New("message not found")
	}

	// In a real app, you might have different rules (e.g., only sender can delete)
	// Here, we allow either participant to delete it for both.
	if msg.SenderID != userID && msg.ReceiverID != userID {
		return errors.New("user not authorized to delete this message")
	}

	if err := s.messageRepo.Delete(messageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	// Broadcast event via WebSocket
	s.hub.BroadcastEvent(&websocket.Event{
		Type:    websocket.EventMessageDeleted,
		Payload: *msg, // Send the old message data to identify it on the client
	})

	return nil
}

func (s *messageService) GetConversation(user1ID, user2ID uint) ([]models.Message, error) {
	return s.messageRepo.GetConversation(user1ID, user2ID)
}
