package services

import (
	"testing"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
)

func TestBuildMessagePushPayload_DataOnlyContract(t *testing.T) {
	createdAt := time.Date(2026, 3, 25, 10, 30, 0, 0, time.UTC)
	payload := models.PushMessagePayload{
		Type:       "chat_message",
		ChatID:     "private-main",
		MessageID:  987,
		SenderID:   1,
		SenderName: "Alexis",
		Content:    "holi",
		CreatedAt:  createdAt,
	}

	msg := buildMessagePushPayload("token-123", payload)

	if msg.Token != "token-123" {
		t.Fatalf("token esperado token-123, se obtuvo %q", msg.Token)
	}

	if msg.Notification != nil {
		t.Fatal("notification debe ser nil para enviar push data-only")
	}

	if msg.Data["type"] != "chat_message" {
		t.Fatalf("type esperado chat_message, se obtuvo %q", msg.Data["type"])
	}
	if msg.Data["chat_id"] != "private-main" {
		t.Fatalf("chat_id esperado private-main, se obtuvo %q", msg.Data["chat_id"])
	}
	if msg.Data["message_id"] != "987" {
		t.Fatalf("message_id esperado 987, se obtuvo %q", msg.Data["message_id"])
	}
	if msg.Data["sender_id"] != "1" {
		t.Fatalf("sender_id esperado 1, se obtuvo %q", msg.Data["sender_id"])
	}
	if msg.Data["sender_name"] != "Alexis" {
		t.Fatalf("sender_name esperado Alexis, se obtuvo %q", msg.Data["sender_name"])
	}
	if msg.Data["content"] != "holi" {
		t.Fatalf("content esperado holi, se obtuvo %q", msg.Data["content"])
	}
	if msg.Data["created_at"] != createdAt.Format(time.RFC3339) {
		t.Fatalf("created_at inesperado: %q", msg.Data["created_at"])
	}

	if msg.Android == nil || msg.Android.Priority != "high" {
		t.Fatal("android.priority debe ser high")
	}

	if msg.APNS == nil {
		t.Fatal("apns no debe ser nil")
	}
	if msg.APNS.Headers["apns-priority"] != "5" {
		t.Fatalf("apns-priority esperado 5, se obtuvo %q", msg.APNS.Headers["apns-priority"])
	}
	if msg.APNS.Headers["apns-push-type"] != "background" {
		t.Fatalf("apns-push-type esperado background, se obtuvo %q", msg.APNS.Headers["apns-push-type"])
	}
	if msg.APNS.Payload == nil || msg.APNS.Payload.Aps == nil || !msg.APNS.Payload.Aps.ContentAvailable {
		t.Fatal("apns payload debe incluir content-available")
	}
}
