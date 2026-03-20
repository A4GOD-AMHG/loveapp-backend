// Paquete websocket implementa la comunicación en tiempo real mediante WebSockets.
package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Constantes de configuración para el manejo de conexiones WebSocket.
const (
	writeWait      = 10 * time.Second                     // Tiempo máximo para escribir un mensaje al cliente
	pongWait       = 60 * time.Second                     // Tiempo máximo para recibir un pong del cliente
	pingPeriod     = (pongWait * 9) / 10                  // Intervalo de envío de pings (debe ser menor a pongWait)
	maxMessageSize = 512                                   // Tamaño máximo permitido de un mensaje entrante (bytes)
)

// upgrader configura la actualización de conexiones HTTP a WebSocket.
// Permite todas las conexiones entrantes (CheckOrigin siempre retorna true).
// En producción se debería implementar una validación de origen más estricta.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Permitir todas las conexiones; en producción validar el origen con una allowlist
		return true
	},
}

// Client representa un cliente WebSocket conectado.
// Actúa como intermediario entre la conexión WebSocket y el Hub central.
type Client struct {
	hub    *Hub            // Referencia al Hub que gestiona este cliente
	conn   *websocket.Conn // Conexión WebSocket subyacente
	send   chan []byte      // Canal con mensajes pendientes de enviar al cliente
	userID uint            // ID del usuario autenticado propietario de esta conexión
}

// readPump escucha mensajes entrantes desde la conexión WebSocket del cliente.
// En esta implementación el cliente no envía mensajes via WebSocket (usa HTTP para eso),
// por lo que este método solo mantiene la conexión viva y detecta desconexiones.
// Al terminar, desregistra al cliente del Hub y cierra la conexión.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Configurar límites y deadline de lectura
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// Renovar el deadline de lectura cada vez que se recibe un pong
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Leer mensajes entrantes solo para mantener la conexión activa
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// writePump envía mensajes desde el canal send del cliente hacia la conexión WebSocket.
// Envía pings periódicos para mantener la conexión activa y detectar clientes desconectados.
// Al terminar (por cierre del canal o error), cierra la conexión WebSocket.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// Establecer deadline de escritura antes de cada mensaje
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// El Hub cerró el canal: enviar mensaje de cierre al cliente
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Escribir el mensaje como texto al cliente WebSocket
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Enviar ping periódico para mantener la conexión activa
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs actualiza una conexión HTTP a WebSocket, registra al cliente en el Hub
// e inicia las goroutines de lectura y escritura para la comunicación bidireccional.
// El userID se usa para identificar al cliente y dirigirle eventos específicos.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID uint) {
	// Actualizar la conexión HTTP a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Crear el cliente y registrarlo en el Hub
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), userID: userID}
	client.hub.register <- client

	// Iniciar goroutines de lectura y escritura de forma concurrente
	// Esto permite que el llamador libere la memoria referenciada
	go client.writePump()
	go client.readPump()
}
