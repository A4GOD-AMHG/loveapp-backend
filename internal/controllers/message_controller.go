// Paquete controllers contiene los manejadores HTTP de la aplicación.
package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
)

// MessageController maneja los endpoints relacionados con mensajes y WebSocket.
type MessageController struct {
	service services.MessageService // Servicio de mensajes con la lógica de negocio
	hub     *websocket.Hub          // Hub de WebSocket para difundir eventos en tiempo real
}

// NewMessageController crea y retorna una nueva instancia de MessageController.
func NewMessageController(service services.MessageService, hub *websocket.Hub) *MessageController {
	return &MessageController{service: service, hub: hub}
}

// ServeWS maneja las solicitudes de conexión WebSocket.
// @Summary Conexión WebSocket
// @Description Actualiza la conexión HTTP a una conexión WebSocket para comunicación en tiempo real.
// @Tags Mensajes
// @Security BearerAuth
// @Success 101 "Cambiando Protocolos"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Router /ws [get]
func (h *MessageController) ServeWS(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	websocket.ServeWs(h.hub, w, r, uint(user.ID))
}

// sendMessageRequest representa el cuerpo de la solicitud para enviar un mensaje.
type sendMessageRequest struct {
	Content string `json:"content"` // Contenido del mensaje a enviar
}

// SendMessage maneja el envío de un nuevo mensaje.
// @Summary Enviar un mensaje
// @Description Envía un mensaje del usuario autenticado al otro usuario del sistema.
// @Tags Mensajes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body sendMessageRequest true "Carga útil del mensaje"
// @Success 201 {object} models.Message "Mensaje enviado exitosamente"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /messages [post]
func (h *MessageController) SendMessage(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	var req sendMessageRequest
	if !decodeJSON(w, r, &req, "Por favor, revisa los datos del mensaje.") {
		return
	}

	msg, err := h.service.SendMessage(uint(user.ID), req.Content)
	switch {
	case err == nil:
		response.JSON(w, http.StatusCreated, msg)
	case errors.Is(err, services.ErrInvalidInput):
		response.BadRequest(w, "Falta el contenido del mensaje.")
	case errors.Is(err, services.ErrNoReceiver):
		response.BadRequest(w, "No se encontró otro usuario para enviar el mensaje.")
	default:
		response.InternalServerError(w, "No pudimos enviar tu mensaje. Inténtalo de nuevo.")
	}
}

// editMessageRequest representa el cuerpo de la solicitud para editar un mensaje.
type editMessageRequest struct {
	Content string `json:"content"` // Nuevo contenido del mensaje
}

// EditMessage maneja la edición de un mensaje existente.
// @Summary Editar un mensaje
// @Description Edita un mensaje enviado por el usuario autenticado, si es dentro de 1 hora de enviado.
// @Tags Mensajes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del mensaje"
// @Param message body editMessageRequest true "Nuevo contenido del mensaje"
// @Success 200 {object} models.Message "Mensaje actualizado exitosamente"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 404 {object} models.ErrorResponse "No encontrado"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /messages/{id} [put]
func (h *MessageController) EditMessage(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "El ID del mensaje no es válido.")
	if !ok {
		return
	}
	var req editMessageRequest
	if !decodeJSON(w, r, &req, "Por favor, revisa el contenido del mensaje.") {
		return
	}

	msg, err := h.service.EditMessage(uint(user.ID), id, req.Content)
	switch {
	case err == nil:
		response.JSON(w, http.StatusOK, msg)
	case errors.Is(err, services.ErrInvalidInput):
		response.BadRequest(w, "El contenido del mensaje no puede estar vacío.")
	case errors.Is(err, services.ErrNotFound):
		response.NotFound(w, "El mensaje que intentas editar no existe.")
	case errors.Is(err, services.ErrForbidden):
		response.Forbidden(w, "No tienes permiso para editar este mensaje.")
	case errors.Is(err, services.ErrExpired):
		response.Forbidden(w, "Ya ha pasado demasiado tiempo para poder editar este mensaje.")
	default:
		response.InternalServerError(w, "No pudimos editar el mensaje. Inténtalo más tarde.")
	}
}

// DeleteMessage maneja la eliminación de un mensaje.
// @Summary Eliminar un mensaje
// @Description Elimina un mensaje para el remitente y el destinatario.
// @Tags Mensajes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del mensaje"
// @Success 204 "Sin Contenido"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 404 {object} models.ErrorResponse "No encontrado"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /messages/{id} [delete]
func (h *MessageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "El ID del mensaje no es válido.")
	if !ok {
		return
	}

	err := h.service.DeleteMessage(uint(user.ID), id)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, services.ErrNotFound):
		response.NotFound(w, "El mensaje que intentas eliminar no existe.")
	case errors.Is(err, services.ErrForbidden):
		response.Forbidden(w, "No tienes permiso para eliminar este mensaje.")
	case errors.Is(err, services.ErrExpired):
		response.Forbidden(w, "Ya ha pasado demasiado tiempo para poder eliminar este mensaje.")
	default:
		response.InternalServerError(w, "No pudimos eliminar el mensaje. Inténtalo más tarde.")
	}
}

// MarkAsRead maneja la acción de marcar un mensaje como leído.
// @Summary Marcar mensaje como leído
// @Description Marca un mensaje como leído por el destinatario.
// @Tags Mensajes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del mensaje"
// @Success 200 {object} map[string]string "Mensaje marcado como leído"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 404 {object} models.ErrorResponse "No encontrado"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /messages/{id}/read [patch]
func (h *MessageController) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	h.updateStatus(w, r, h.service.MarkAsRead, "leído")
}

// MarkAsDelivered maneja la acción de marcar un mensaje como entregado.
// @Summary Marcar mensaje como entregado
// @Description Marca un mensaje como entregado al destinatario.
// @Tags Mensajes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del mensaje"
// @Success 200 {object} map[string]string "Mensaje marcado como entregado"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 404 {object} models.ErrorResponse "No encontrado"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /messages/{id}/delivered [patch]
func (h *MessageController) MarkAsDelivered(w http.ResponseWriter, r *http.Request) {
	h.updateStatus(w, r, h.service.MarkAsDelivered, "entregado")
}

// updateStatus comparte el flujo de MarkAsRead y MarkAsDelivered.
func (h *MessageController) updateStatus(w http.ResponseWriter, r *http.Request, mark func(uint, int64) error, label string) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "El ID del mensaje no es válido.")
	if !ok {
		return
	}

	err := mark(uint(user.ID), id)
	switch {
	case err == nil:
		response.JSON(w, http.StatusOK, map[string]string{"message": "Mensaje marcado como " + label})
	case errors.Is(err, services.ErrNotFound):
		response.NotFound(w, "El mensaje no existe.")
	case errors.Is(err, services.ErrForbidden):
		response.Forbidden(w, "No tienes permiso para marcar este mensaje como "+label+".")
	default:
		response.InternalServerError(w, "No pudimos marcar el mensaje como "+label+".")
	}
}

// GetConversation maneja la recuperación del historial de mensajes entre dos usuarios.
// @Summary Obtener conversación
// @Description Recupera el historial de mensajes entre el usuario autenticado y el otro usuario del sistema, con paginación.
// @Tags Mensajes
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página (por defecto 1)"
// @Param per_page query int false "Mensajes por página (por defecto 10)"
// @Success 200 {array} models.Message "Historial de conversación"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /messages/conversation [get]
func (h *MessageController) GetConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	// Si los parámetros son inválidos, el repositorio aplica los valores predeterminados.
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	messages, err := h.service.GetConversation(uint(user.ID), page, perPage)
	switch {
	case err == nil:
		if messages == nil {
			messages = []models.Message{}
		}
		response.JSON(w, http.StatusOK, messages)
	case errors.Is(err, services.ErrNoReceiver):
		response.JSON(w, http.StatusOK, []models.Message{})
	default:
		response.InternalServerError(w, "No pudimos cargar la conversación. Inténtalo más tarde.")
	}
}

// GetUnreadCount retorna la cantidad total de mensajes no leídos del usuario autenticado.
func (h *MessageController) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	unreadCount, err := h.service.GetUnreadCount(uint(user.ID))
	if err != nil {
		response.InternalServerError(w, "No pudimos calcular tus mensajes no leídos.")
		return
	}
	response.JSON(w, http.StatusOK, models.UnreadCountResponse{UnreadCount: unreadCount})
}
