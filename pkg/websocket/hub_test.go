package websocket

import (
	"testing"
	"time"
)

func newTestClient(hub *Hub, userID uint) *Client {
	c := &Client{hub: hub, send: make(chan []byte, 4), userID: userID}
	hub.register <- c
	return c
}

func receive(t *testing.T, c *Client) {
	t.Helper()
	select {
	case <-c.send:
	case <-time.After(time.Second):
		t.Fatalf("cliente %d no recibió el evento", c.userID)
	}
}

// Una reconexión no debe perder eventos cuando la conexión vieja se cierra después.
func TestHub_ReconnectKeepsNewClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	old := newTestClient(hub, 1)
	fresh := newTestClient(hub, 1)
	hub.unregister <- old

	hub.BroadcastEvent(&Event{Type: EventMessageSent, SenderID: 2, ReceiverID: 1})
	receive(t, fresh)
	if _, open := <-old.send; open {
		t.Fatal("el canal del cliente viejo debería estar cerrado")
	}
}

func TestHub_BroadcastReachesSenderAndReceiver(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	sender := newTestClient(hub, 1)
	receiver := newTestClient(hub, 2)
	other := newTestClient(hub, 3)

	hub.BroadcastEvent(&Event{Type: EventMessageSent, SenderID: 1, ReceiverID: 2})
	receive(t, sender)
	receive(t, receiver)
	select {
	case <-other.send:
		t.Fatal("un usuario ajeno no debe recibir el evento")
	case <-time.After(50 * time.Millisecond):
	}
}
