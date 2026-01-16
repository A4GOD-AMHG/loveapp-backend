package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[uint]*Client
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
	// Mutex to protect clients map
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Printf("Client connected: %d", client.userID)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				log.Printf("Client disconnected: %d", client.userID)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			var event Event
			if err := json.Unmarshal(message, &event); err != nil {
				log.Printf("error unmarshalling broadcast message: %v", err)
				continue
			}

			h.mu.RLock()
			// Attempt to send to receiver
			if client, ok := h.clients[event.Payload.ReceiverID]; ok {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.userID)
				}
			}
			// Also send to sender to confirm message was sent and update their UI
			if client, ok := h.clients[event.Payload.SenderID]; ok {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.userID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastEvent sends an event to the relevant clients.
func (h *Hub) BroadcastEvent(event *Event) {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("error marshalling event: %v", err)
		return
	}
	h.broadcast <- bytes
}
