package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"google.golang.org/api/option"
)

// PushService abstrae el envío de notificaciones push.
type PushService interface {
	SendNewMessage(tokens []models.DevicePushToken, payload models.PushMessagePayload) error
}

type pushService struct {
	client *messaging.Client
}

// NewPushService crea un servicio de push usando Firebase Admin SDK.
func NewPushService() PushService {
	credentialsFile := config.AppConfig.Push.CredentialsFile
	if credentialsFile == "" {
		log.Printf("push notifications omitidas: FIREBASE_CREDENTIALS_FILE no configurado")
		return &pushService{}
	}

	if _, err := os.Stat(credentialsFile); err != nil {
		log.Printf("push notifications omitidas: no se encontró el archivo de credenciales %q: %v", credentialsFile, err)
		return &pushService{}
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Printf("push notifications omitidas: error inicializando Firebase: %v", err)
		return &pushService{}
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("push notifications omitidas: error creando cliente de Messaging: %v", err)
		return &pushService{}
	}

	return &pushService{client: client}
}

func (s *pushService) SendNewMessage(tokens []models.DevicePushToken, payload models.PushMessagePayload) error {
	if len(tokens) == 0 || s.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	messages := make([]*messaging.Message, 0, len(tokens))
	for _, token := range tokens {
		messages = append(messages, &messaging.Message{
			Token: token.PushToken,
			Data: map[string]string{
				"type":        payload.Type,
				"message_id":  strconv.FormatUint(uint64(payload.MessageID), 10),
				"sender_id":   strconv.FormatUint(uint64(payload.SenderID), 10),
				"sender_name": payload.SenderName,
				"content":     payload.Content,
				"created_at":  payload.CreatedAt.Format(time.RFC3339),
			},
			Notification: &messaging.Notification{
				Title: payload.SenderName,
				Body:  payload.Content,
			},
		})
	}

	batchResponse, err := s.client.SendEach(ctx, messages)
	if err != nil {
		return fmt.Errorf("failed to send push notification batch: %w", err)
	}

	for i, resp := range batchResponse.Responses {
		if resp.Success {
			continue
		}
		log.Printf("error enviando push a token %q: %v", tokens[i].PushToken, resp.Error)
	}

	return nil
}
