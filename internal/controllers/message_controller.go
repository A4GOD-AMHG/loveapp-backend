package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
	"github.com/gorilla/mux"
)

type MessageController struct {
	service services.MessageService
	hub     *websocket.Hub
}

func NewMessageController(service services.MessageService, hub *websocket.Hub) *MessageController {
	return &MessageController{
		service: service,
		hub:     hub,
	}
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
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Tu sesión ha expirado, por favor inicia sesión de nuevo.")
		return
	}
	websocket.ServeWs(h.hub, w, r, uint(user.ID))
}

type sendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

// SendMessage maneja el envío de un nuevo mensaje.
// @Summary Enviar un mensaje
// @Description Envía un mensaje del usuario autenticado a otro usuario.
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
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Tu sesión ha expirado, por favor inicia sesión de nuevo.")
		return
	}

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos del mensaje.")
		return
	}

	if req.Content == "" || req.ReceiverID == 0 {
		response.BadRequest(w, "Falta el destinatario o el contenido del mensaje.")
		return
	}

	msg, err := h.service.SendMessage(uint(user.ID), req.ReceiverID, req.Content)
	if err != nil {
		response.InternalServerError(w, "No pudimos enviar tu mensaje. Inténtalo de nuevo.")
		return
	}

	response.JSON(w, http.StatusCreated, msg)
}

type editMessageRequest struct {
	Content string `json:"content"`
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
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Tu sesión ha expirado, por favor inicia sesión de nuevo.")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.BadRequest(w, "El ID del mensaje no es válido.")
		return
	}

	var req editMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa el contenido del mensaje.")
		return
	}

	if req.Content == "" {
		response.BadRequest(w, "El contenido del mensaje no puede estar vacío.")
		return
	}

	msg, err := h.service.EditMessage(uint(user.ID), id, req.Content)
	if err != nil {
		if err.Error() == "message not found" {
			response.NotFound(w, "El mensaje que intentas editar no existe.")
		} else if err.Error() == "user not authorized to edit this message" {
			response.Forbidden(w, "No tienes permiso para editar este mensaje.")
		} else if err.Error() == "message can no longer be edited" {
			response.Forbidden(w, "Ya ha pasado demasiado tiempo para poder editar este mensaje.")
		} else {
			response.InternalServerError(w, "No pudimos editar el mensaje. Inténtalo más tarde.")
		}
		return
	}

	response.JSON(w, http.StatusOK, msg)
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
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Tu sesión ha expirado, por favor inicia sesión de nuevo.")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.BadRequest(w, "El ID del mensaje no es válido.")
		return
	}

	err = h.service.DeleteMessage(uint(user.ID), id)
	if err != nil {
		if err.Error() == "message not found" {
			response.NotFound(w, "El mensaje que intentas eliminar no existe.")
		} else if err.Error() == "user not authorized to delete this message" {
			response.Forbidden(w, "No tienes permiso para eliminar este mensaje.")
		} else {
			response.InternalServerError(w, "No pudimos eliminar el mensaje. Inténtalo más tarde.")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetConversation maneja la recuperación del historial de mensajes entre dos usuarios.
// @Summary Obtener conversación
// @Description Recupera el historial completo de mensajes entre el usuario autenticado y otro usuario.
// @Tags Mensajes
// @Produce json
// @Security BearerAuth
// @Param userId path int true "El ID del otro usuario"
// @Success 200 {array} models.Message "Historial de conversación"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /messages/conversation/{userId} [get]
func (h *MessageController) GetConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Tu sesión ha expirado, por favor inicia sesión de nuevo.")
		return
	}

	vars := mux.Vars(r)
	otherUserID, err := strconv.ParseUint(vars["userId"], 10, 32)
	if err != nil {
		response.BadRequest(w, "El ID de usuario no es válido.")
		return
	}

	messages, err := h.service.GetConversation(uint(user.ID), uint(otherUserID))
	if err != nil {
		response.InternalServerError(w, "No pudimos cargar la conversación. Inténtalo más tarde.")
		return
	}

	response.JSON(w, http.StatusOK, messages)
}
