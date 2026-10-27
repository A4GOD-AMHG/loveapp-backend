package services

import (
	"errors"
	"testing"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/testhelpers"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
)

func TestNormalizeSenderName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "alexis lowercase", input: "alexis", expected: "Alexis"},
		{name: "alexis with suffix", input: "alexis :)", expected: "Alexis"},
		{name: "alexis mixed case", input: "ALeXiS", expected: "Alexis"},
		{name: "anyel lowercase", input: "anyel", expected: "Anyel"},
		{name: "anyel with spaces", input: "  anyel  ", expected: "Anyel"},
		{name: "anyel mixed case", input: "AnYeL", expected: "Anyel"},
		{name: "unknown user", input: "Carlos", expected: "Carlos"},
		{name: "unknown trimmed", input: "  Laura  ", expected: "Laura"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeSenderName(tt.input)
			if got != tt.expected {
				t.Fatalf("resultado inesperado: esperado %q, obtenido %q", tt.expected, got)
			}
		})
	}
}

type fakeHub struct{ events []*websocket.Event }

func (h *fakeHub) BroadcastEvent(e *websocket.Event) { h.events = append(h.events, e) }

type fakePush struct {
	calls  chan models.PushMessagePayload
	tokens chan []models.DevicePushToken
}

func (p *fakePush) SendNewMessage(tokens []models.DevicePushToken, payload models.PushMessagePayload) error {
	p.tokens <- tokens
	p.calls <- payload
	return nil
}

func TestSendMessage_NotifiesReceiverDevices(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()
	anyelID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "x")
	alexisID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "x")

	deviceRepo := repository.NewDevicePushTokenRepository()
	_ = deviceRepo.Upsert(&models.DevicePushToken{UserID: alexisID, Platform: "android", PushToken: "tok-alexis", DeviceName: "Pixel"})
	_ = deviceRepo.Upsert(&models.DevicePushToken{UserID: anyelID, Platform: "android", PushToken: "tok-anyel", DeviceName: "Moto"})

	hub := &fakeHub{}
	push := &fakePush{calls: make(chan models.PushMessagePayload, 1), tokens: make(chan []models.DevicePushToken, 1)}
	svc := NewMessageService(repository.NewMessageRepository(), repository.NewUserRepository(), deviceRepo, push, hub)

	msg, err := svc.SendMessage(uint(anyelID), "te quiero")
	if err != nil {
		t.Fatal(err)
	}
	if len(hub.events) != 1 || hub.events[0].ReceiverID != uint(alexisID) {
		t.Fatalf("debe difundir un evento hacia Alexis: %+v", hub.events)
	}

	select {
	case tokens := <-push.tokens:
		if len(tokens) != 1 || tokens[0].PushToken != "tok-alexis" {
			t.Fatalf("el push debe ir solo a los dispositivos de Alexis: %+v", tokens)
		}
		payload := <-push.calls
		if payload.SenderName != "Anyel" || payload.Content != "te quiero" || payload.MessageID != msg.ID {
			t.Fatalf("payload push inesperado: %+v", payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no se envió la notificación push")
	}

	if _, err := svc.SendMessage(uint(anyelID), "  "); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("mensaje vacío debe ser ErrInvalidInput, fue %v", err)
	}
	if err := svc.MarkAsRead(uint(anyelID), int64(msg.ID)); !errors.Is(err, ErrForbidden) {
		t.Fatalf("el remitente no puede marcar leído: %v", err)
	}
	if err := svc.MarkAsRead(uint(alexisID), int64(msg.ID)); err != nil {
		t.Fatal(err)
	}
	// read -> delivered no debe retroceder el estado ni emitir evento
	before := len(hub.events)
	if err := svc.MarkAsDelivered(uint(alexisID), int64(msg.ID)); err != nil || len(hub.events) != before {
		t.Fatalf("delivered tras read no debe emitir evento: %v", err)
	}
}
