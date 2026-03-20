// Paquete controllers contiene los manejadores HTTP de la aplicación.
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
	"github.com/gorilla/mux"
)

// TodoController maneja los endpoints relacionados con la gestión de tareas.
type TodoController struct {
	todoService *services.TodoService // Servicio de tareas con la lógica de negocio
}

// NewTodoController crea y retorna una nueva instancia de TodoController
// con su servicio de tareas ya inicializado.
func NewTodoController() *TodoController {
	return &TodoController{
		todoService: services.NewTodoService(),
	}
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
	// Obtener el usuario autenticado desde el contexto (inyectado por el middleware)
	user := r.Context().Value("user").(*models.User)

	var req models.CreateTodoRequest

	// Decodificar el cuerpo JSON con los datos de la nueva tarea
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos ingresados.")
		return
	}

	// Delegar la creación al servicio con validación de título
	createResponse, err := h.todoService.CreateTodo(user.ID, user.Username, &req)
	if err != nil {
		if err.Error() == "el título es requerido" {
			response.BadRequest(w, "El título de la tarea no puede estar vacío.")
			return
		}
		response.InternalServerError(w, "Tuvimos un problema al crear la tarea. Inténtalo más tarde.")
		return
	}

	response.JSON(w, http.StatusCreated, createResponse)
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
	// Obtener el usuario autenticado desde el contexto
	user := r.Context().Value("user").(*models.User)

	// Leer los filtros y parámetros de paginación desde la query string
	status := r.URL.Query().Get("status")
	creatorID := r.URL.Query().Get("creator_id")
	sortOrder := r.URL.Query().Get("sort_order")
	search := r.URL.Query().Get("search")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// Delegar la consulta filtrada al servicio de tareas
	listResponse, err := h.todoService.GetTodos(status, creatorID, user.Username, sortOrder, search, page, limit)
	if err != nil {
		if err.Error() == "ID de creador inválido" {
			response.BadRequest(w, "El ID del creador no es válido.")
			return
		}
		response.InternalServerError(w, "No pudimos cargar las tareas. Inténtalo de nuevo.")
		return
	}

	response.JSON(w, http.StatusOK, listResponse)
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
	// Obtener el ID del usuario autenticado desde el contexto
	userID := r.Context().Value("user_id").(int64)

	// Parsear el ID de la tarea desde los parámetros de la URL
	vars := mux.Vars(r)
	todoIDStr := vars["id"]
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "El ID de la tarea no es válido.")
		return
	}

	// Decodificar el cuerpo JSON con los nuevos datos de la tarea
	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos ingresados.")
		return
	}

	// Validar que el título no esté vacío
	if req.Title == "" {
		response.BadRequest(w, "El título de la tarea no puede estar vacío.")
		return
	}

	// Delegar la actualización al servicio (solo el creador puede editar)
	updatedTodo, err := h.todoService.UpdateTodo(todoID, userID, &req)
	if err != nil {
		if err.Error() == "todo no encontrado" {
			response.NotFound(w, "No encontramos la tarea que quieres editar.")
			return
		}
		if err.Error() == "solo el creador puede editar este todo" {
			response.Forbidden(w, "Solo quien creó la tarea puede editarla.")
			return
		}
		response.InternalServerError(w, "No pudimos editar la tarea. Inténtalo más tarde.")
		return
	}

	response.JSON(w, http.StatusOK, updatedTodo)
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
	// Obtener el nombre de usuario autenticado desde el contexto
	username := r.Context().Value("username").(string)

	// Parsear el ID de la tarea desde los parámetros de la URL
	vars := mux.Vars(r)
	todoIDStr := vars["id"]
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "El ID de la tarea no es válido.")
		return
	}

	// Decodificar el nuevo estado (completado o no) desde el cuerpo JSON
	var req models.UpdateTodoStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos ingresados.")
		return
	}

	// Delegar la actualización de estado al servicio con validaciones de usuario
	statusResponse, err := h.todoService.UpdateTodoStatus(todoID, username, req.Completed)
	if err != nil {
		if err.Error() == "todo no encontrado" {
			response.NotFound(w, "No encontramos la tarea que buscas.")
			return
		}
		if err.Error() == "usuario no autorizado para marcar todos" {
			response.Forbidden(w, "No tienes permiso para actualizar esta tarea.")
			return
		}
		if err.Error() == "la tarea ya está completada por ambos y no se puede desmarcar" {
			response.BadRequest(w, "¡La tarea ya está completada por ambos! No se puede desmarcar. 🎉")
			return
		}
		response.InternalServerError(w, "Hubo un error al actualizar la tarea. Inténtalo de nuevo.")
		return
	}

	response.JSON(w, http.StatusOK, statusResponse)
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
	// Obtener el ID del usuario autenticado desde el contexto
	userID := r.Context().Value("user_id").(int64)

	// Parsear el ID de la tarea desde los parámetros de la URL
	vars := mux.Vars(r)
	todoIDStr := vars["id"]
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "El ID de la tarea no es válido.")
		return
	}

	// Delegar la eliminación al servicio (solo el creador puede eliminar)
	err = h.todoService.DeleteTodo(todoID, userID)
	if err != nil {
		if err.Error() == "todo no encontrado" {
			response.NotFound(w, "No encontramos la tarea que quieres eliminar.")
			return
		}
		if err.Error() == "solo el creador puede eliminar este todo" {
			response.Forbidden(w, "Solo quien creó la tarea puede eliminarla.")
			return
		}
		response.InternalServerError(w, "No pudimos eliminar la tarea. Inténtalo más tarde.")
		return
	}

	response.Success(w, "¡Tarea eliminada con éxito! 👋", nil)
}
