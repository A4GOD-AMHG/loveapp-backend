// Paquete websocket implementa la comunicación en tiempo real mediante WebSockets.
package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub centraliza la gestión de clientes WebSocket conectados y la distribución de eventos.
// Mantiene un mapa de clientes activos indexado por userID para envíos dirigidos.
type Hub struct {
	clients    map[uint]*Client // Mapa de clientes conectados: userID → *Client
	broadcast  chan *Event      // Canal de eventos a distribuir a los clientes
	register   chan *Client     // Canal para registrar nuevos clientes al conectarse
	unregister chan *Client     // Canal para eliminar clientes al desconectarse
	mu         sync.RWMutex     // Mutex de lectura/escritura para proteger el mapa de clientes
}

// NewHub crea e inicializa un nuevo Hub con todos sus canales y el mapa de clientes vacío.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Event),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint]*Client),
	}
}

// Run inicia el bucle principal del Hub en una goroutine.
// Escucha continuamente tres tipos de eventos:
//   - register: registra un nuevo cliente y lo agrega al mapa
//   - unregister: elimina un cliente desconectado y cierra su canal de envío
//   - broadcast: distribuye un mensaje JSON al remitente y al destinatario del evento
//
// Debe ejecutarse en una goroutine separada: go hub.Run()
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Registrar el nuevo cliente en el mapa protegido por mutex
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Printf("Client connected: %d", client.userID)

		case client := <-h.unregister:
			// Eliminar el cliente del mapa y cerrar su canal de envío si estaba registrado
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				log.Printf("Client disconnected: %d", client.userID)
			}
			h.mu.Unlock()

		case event := <-h.broadcast:
			message, err := json.Marshal(event)
			if err != nil {
				log.Printf("error marshalling broadcast event: %v", err)
				continue
			}

			h.mu.Lock()

			// Enviar el evento al destinatario si está conectado
			if client, ok := h.clients[event.ReceiverID]; ok {
				select {
				case client.send <- message:
				default:
					// Canal lleno: desconectar al cliente para evitar bloqueos
					close(client.send)
					delete(h.clients, client.userID)
				}
			}

			// También enviar al remitente para confirmar el envío y actualizar su UI
			if client, ok := h.clients[event.SenderID]; ok {
				select {
				case client.send <- message:
				default:
					// Canal lleno: desconectar al cliente para evitar bloqueos
					close(client.send)
					delete(h.clients, client.userID)
				}
			}

			h.mu.Unlock()
		}
	}
}

// BroadcastEvent encola un evento para que el hub lo serialice y distribuya.
// Los clientes relevantes (remitente y destinatario) recibirán el evento en tiempo real.
func (h *Hub) BroadcastEvent(event *Event) {
	h.broadcast <- event
}
