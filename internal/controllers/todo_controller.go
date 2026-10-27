// Paquete controllers contiene los manejadores HTTP de la aplicación.
package controllers

import (
	"errors"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// TodoController maneja los endpoints relacionados con la gestión de tareas.
type TodoController struct {
	todoService *services.TodoService // Servicio de tareas con la lógica de negocio
}

// NewTodoController crea y retorna una nueva instancia de TodoController
// con su servicio de tareas ya inicializado.
func NewTodoController(todoService *services.TodoService) *TodoController {
	return &TodoController{todoService: todoService}
}

// CreateTodo maneja la creación de tareas
// @Summary Crear una nueva tarea
// @Description Crear una nueva tarea (solo para usuarios autorizados: anyel, alexis)
// @Tags Tareas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param todo body models.CreateTodoRequest true "Solicitud de creación de tarea"
// @Success 201 {object} models.CreateTodoResponse "Tarea creada exitosamente"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /todos [post]
func (h *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	var req models.CreateTodoRequest
	if !decodeJSON(w, r, &req, "Por favor, revisa los datos ingresados.") {
		return
	}

	createResponse, err := h.todoService.CreateTodo(user.ID, &req)
	switch {
	case err == nil:
		response.JSON(w, http.StatusCreated, createResponse)
	case errors.Is(err, services.ErrInvalidInput):
		response.BadRequest(w, "El título de la tarea no puede estar vacío.")
	default:
		response.InternalServerError(w, "Tuvimos un problema al crear la tarea. Inténtalo más tarde.")
	}
}

// GetTodos maneja el listado de tareas con filtros y ordenamiento
// @Summary Listar tareas
// @Description Obtener una lista de tareas con filtros opcionales por estado, creador y ordenamiento por fecha de creación.
// @Tags Tareas
// @Produce json
// @Security BearerAuth
// @Param creator_id query int false "Filtrar por ID de usuario creador."
// @Param status query string false "Filtrar por estado." Enums(all, completed, incompleted, completed_by_me) default(all)
// @Param search query string false "Buscar por título o descripción."
// @Param sort_order query string false "Orden por fecha de creación." Enums(asc, desc) default(desc)
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} models.TodoListResponse "Tareas recuperadas exitosamente"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /todos [get]
func (h *TodoController) GetTodos(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	q := r.URL.Query()

	listResponse, err := h.todoService.GetTodos(
		q.Get("status"), q.Get("creator_id"), user.Username,
		q.Get("sort_order"), q.Get("search"), q.Get("page"), q.Get("limit"),
	)
	switch {
	case err == nil:
		response.JSON(w, http.StatusOK, listResponse)
	case errors.Is(err, services.ErrInvalidInput):
		response.BadRequest(w, "El ID del creador no es válido.")
	default:
		response.InternalServerError(w, "No pudimos cargar las tareas. Inténtalo de nuevo.")
	}
}

// UpdateTodo maneja la actualización de tareas (título/descripción)
// @Summary Actualizar detalles de la tarea
// @Description Actualizar título y descripción de la tarea (solo el creador puede actualizar sus tareas)
// @Tags Tareas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la tarea"
// @Param todo body models.UpdateTodoRequest true "Solicitud de actualización de tarea"
// @Success 200 {object} models.Todo "Tarea actualizada exitosamente"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 404 {object} models.ErrorResponse "Tarea no encontrada"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /todos/{id} [put]
func (h *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	todoID, ok := pathID(w, r, "El ID de la tarea no es válido.")
	if !ok {
		return
	}
	var req models.UpdateTodoRequest
	if !decodeJSON(w, r, &req, "Por favor, revisa los datos ingresados.") {
		return
	}

	updatedTodo, err := h.todoService.UpdateTodo(todoID, user.ID, &req)
	switch {
	case err == nil:
		response.JSON(w, http.StatusOK, updatedTodo)
	case errors.Is(err, services.ErrInvalidInput):
		response.BadRequest(w, "El título de la tarea no puede estar vacío.")
	case errors.Is(err, services.ErrNotFound):
		response.NotFound(w, "No encontramos la tarea que quieres editar.")
	case errors.Is(err, services.ErrForbidden):
		response.Forbidden(w, "Solo quien creó la tarea puede editarla.")
	default:
		response.InternalServerError(w, "No pudimos editar la tarea. Inténtalo más tarde.")
	}
}

// UpdateTodoStatus maneja la actualización del estado de la tarea
// @Summary Actualizar estado de la tarea
// @Description Actualizar el estado de la tarea (completado) para el usuario autenticado
// @Tags Tareas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la tarea"
// @Param status body models.UpdateTodoStatusRequest true "Solicitud de actualización de estado de tarea"
// @Success 200 {object} models.CompleteTodoResponse "Estado de la tarea actualizado"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 404 {object} models.ErrorResponse "Tarea no encontrada"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /todos/{id} [patch]
func (h *TodoController) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	todoID, ok := pathID(w, r, "El ID de la tarea no es válido.")
	if !ok {
		return
	}
	var req models.UpdateTodoStatusRequest
	if !decodeJSON(w, r, &req, "Por favor, revisa los datos ingresados.") {
		return
	}

	statusResponse, err := h.todoService.UpdateTodoStatus(todoID, user.Username, req.Completed)
	switch {
	case err == nil:
		response.JSON(w, http.StatusOK, statusResponse)
	case errors.Is(err, services.ErrNotFound):
		response.NotFound(w, "No encontramos la tarea que buscas.")
	case errors.Is(err, services.ErrForbidden):
		response.Forbidden(w, "No tienes permiso para actualizar esta tarea.")
	case errors.Is(err, services.ErrAlreadyCompleted):
		response.BadRequest(w, "¡La tarea ya está completada por ambos! No se puede desmarcar. 🎉")
	default:
		response.InternalServerError(w, "Hubo un error al actualizar la tarea. Inténtalo de nuevo.")
	}
}

// DeleteTodo maneja la eliminación de tareas
// @Summary Eliminar una tarea
// @Description Eliminar una tarea (solo el creador puede eliminar sus tareas)
// @Tags Tareas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la tarea"
// @Success 200 {object} models.SuccessResponse "Tarea eliminada exitosamente"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 404 {object} models.ErrorResponse "Tarea no encontrada"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /todos/{id} [delete]
func (h *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}
	todoID, ok := pathID(w, r, "El ID de la tarea no es válido.")
	if !ok {
		return
	}

	err := h.todoService.DeleteTodo(todoID, user.ID)
	switch {
	case err == nil:
		response.Success(w, "¡Tarea eliminada con éxito! 👋", nil)
	case errors.Is(err, services.ErrNotFound):
		response.NotFound(w, "No encontramos la tarea que quieres eliminar.")
	case errors.Is(err, services.ErrForbidden):
		response.Forbidden(w, "Solo quien creó la tarea puede eliminarla.")
	default:
		response.InternalServerError(w, "No pudimos eliminar la tarea. Inténtalo más tarde.")
	}
}
