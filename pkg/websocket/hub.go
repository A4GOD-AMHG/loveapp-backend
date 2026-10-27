// Paquete websocket implementa la comunicación en tiempo real mediante WebSockets.
package websocket

import (
	"encoding/json"
	"log"
)

// Hub centraliza los clientes WebSocket conectados y la distribución de eventos.
// Un usuario puede tener varias conexiones (varios dispositivos o una reconexión
// que llega antes de que el servidor detecte la caída de la anterior).
// Todo el estado lo modifica únicamente la goroutine de Run, por eso no usa mutex.
type Hub struct {
	clients    map[uint]map[*Client]struct{}
	broadcast  chan *Event
	register   chan *Client
	unregister chan *Client
}

// NewHub crea un Hub vacío.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Event),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint]map[*Client]struct{}),
	}
}

// Run procesa registros, desregistros y eventos. Debe ejecutarse con: go hub.Run()
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]struct{})
			}
			h.clients[client.userID][client] = struct{}{}
			log.Printf("Client connected: %d", client.userID)

		case client := <-h.unregister:
			if h.remove(client) {
				log.Printf("Client disconnected: %d", client.userID)
			}

		case event := <-h.broadcast:
			message, err := json.Marshal(event)
			if err != nil {
				log.Printf("error marshalling broadcast event: %v", err)
				continue
			}
			h.sendTo(event.ReceiverID, message)
			if event.SenderID != event.ReceiverID {
				h.sendTo(event.SenderID, message)
			}
		}
	}
}

func (h *Hub) sendTo(userID uint, message []byte) {
	for client := range h.clients[userID] {
		select {
		case client.send <- message:
		default:
			// Buffer lleno: el cliente está atascado, se desconecta.
			h.remove(client)
		}
	}
}

// remove elimina al cliente y cierra su canal una sola vez.
func (h *Hub) remove(client *Client) bool {
	conns := h.clients[client.userID]
	if _, ok := conns[client]; !ok {
		return false
	}
	delete(conns, client)
	if len(conns) == 0 {
		delete(h.clients, client.userID)
	}
	close(client.send)
	return true
}

// BroadcastEvent encola un evento para el remitente y el destinatario.
func (h *Hub) BroadcastEvent(event *Event) {
	h.broadcast <- event
}
