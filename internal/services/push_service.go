package services

import (
	"context"
	"encoding/json"
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

type firebaseServiceAccount struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain,omitempty"`
}

// NewPushService crea un servicio de push usando Firebase Admin SDK.
func NewPushService() PushService {
	ctx := context.Background()
	pushConfig := config.AppConfig.Push
	appConfig := &firebase.Config{
		ProjectID: pushConfig.ProjectID,
	}

	var app *firebase.App
	var err error

	if hasFirebaseEnvCredentials(pushConfig) {
		credentialsJSON, marshalErr := json.Marshal(firebaseServiceAccount{
			Type:                    pushConfig.Type,
			ProjectID:               pushConfig.ProjectID,
			PrivateKeyID:            pushConfig.PrivateKeyID,
			PrivateKey:              pushConfig.PrivateKey,
			ClientEmail:             pushConfig.ClientEmail,
			ClientID:                pushConfig.ClientID,
			AuthURI:                 pushConfig.AuthURI,
			TokenURI:                pushConfig.TokenURI,
			AuthProviderX509CertURL: pushConfig.AuthProviderX509CertURL,
			ClientX509CertURL:       pushConfig.ClientX509CertURL,
			UniverseDomain:          pushConfig.UniverseDomain,
		})
		if marshalErr != nil {
			log.Printf("push notifications omitidas: error serializando credenciales Firebase: %v", marshalErr)
			return &pushService{}
		}

		app, err = firebase.NewApp(ctx, appConfig, option.WithCredentialsJSON(credentialsJSON))
	} else {
		credentialsFile := pushConfig.CredentialsFile
		if credentialsFile == "" {
			log.Printf("push notifications omitidas: credenciales Firebase no configuradas")
			return &pushService{}
		}

		if _, statErr := os.Stat(credentialsFile); statErr != nil {
			log.Printf("push notifications omitidas: no se encontró el archivo de credenciales %q: %v", credentialsFile, statErr)
			return &pushService{}
		}

		app, err = firebase.NewApp(ctx, appConfig, option.WithCredentialsFile(credentialsFile))
	}

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

func hasFirebaseEnvCredentials(cfg config.PushConfig) bool {
	return cfg.Type != "" &&
		cfg.ProjectID != "" &&
		cfg.PrivateKeyID != "" &&
		cfg.PrivateKey != "" &&
		cfg.ClientEmail != "" &&
		cfg.ClientID != "" &&
		cfg.AuthURI != "" &&
		cfg.TokenURI != "" &&
		cfg.AuthProviderX509CertURL != "" &&
		cfg.ClientX509CertURL != ""
}

func (s *pushService) SendNewMessage(tokens []models.DevicePushToken, payload models.PushMessagePayload) error {
	if len(tokens) == 0 || s.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	messages := make([]*messaging.Message, 0, len(tokens))
	for _, token := range tokens {
		messages = append(messages, buildMessagePushPayload(token.PushToken, payload))
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

func buildMessagePushPayload(token string, payload models.PushMessagePayload) *messaging.Message {
	return &messaging.Message{
		Token: token,
		Data: map[string]string{
			"type":        payload.Type,
			"chat_id":     payload.ChatID,
			"message_id":  strconv.FormatUint(uint64(payload.MessageID), 10),
			"sender_id":   strconv.FormatUint(uint64(payload.SenderID), 10),
			"sender_name": payload.SenderName,
			"content":     payload.Content,
			"created_at":  payload.CreatedAt.Format(time.RFC3339),
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority":  "5",
				"apns-push-type": "background",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
			},
		},
	}
}
