package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
	"github.com/gorilla/mux"
)

// TodoHandler handles todo endpoints
type TodoHandler struct {
	todoService *services.TodoService
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		todoService: services.NewTodoService(),
	}
}

// CreateTodo handles todo creation
// @Summary Create a new todo
// @Description Create a new todo item (only for authorized users: anyel, alexis)
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param todo body models.CreateTodoRequest true "Todo creation request"
// @Success 201 {object} models.CreateTodoResponse "Todo created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}
	
	var req models.CreateTodoRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Datos de entrada inválidos")
		return
	}
	
	// Create todo
	createResponse, err := h.todoService.CreateTodo(user.ID, user.Username, &req)
	if err != nil {
		if err.Error() == "usuario no autorizado para crear todos" {
			response.Forbidden(w, err.Error())
			return
		}
		if err.Error() == "el título es requerido" {
			response.BadRequest(w, err.Error())
			return
		}
		response.InternalServerError(w, err.Error())
		return
	}
	
	response.JSON(w, http.StatusCreated, createResponse)
}

// ListTodos handles todo listing with filters
// @Summary List todos
// @Description Get a list of todos with optional filters
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status" Enums(all, pending, completed) default(all)
// @Param creator_id query int false "Filter by creator user ID"
// @Success 200 {object} models.TodoListResponse "Todos retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /todos [get]
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	status := r.URL.Query().Get("status")
	creatorID := r.URL.Query().Get("creator_id")
	
	// List todos
	listResponse, err := h.todoService.ListTodos(status, creatorID)
	if err != nil {
		if err.Error() == "ID de creador inválido" {
			response.BadRequest(w, err.Error())
			return
		}
		response.InternalServerError(w, err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, listResponse)
}

// CompleteTodo handles marking a todo as completed
// @Summary Mark todo as completed
// @Description Mark a todo as completed by the authenticated user (both users must mark it to be fully completed)
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} models.CompleteTodoResponse "Todo marked as completed"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Todo not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /todos/{id}/complete [post]
func (h *TodoHandler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	username, ok := r.Context().Value("username").(string)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}
	
	// Get todo ID from URL
	vars := mux.Vars(r)
	todoIDStr := vars["id"]
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "ID de todo inválido")
		return
	}
	
	// Complete todo
	completeResponse, err := h.todoService.CompleteTodo(todoID, username)
	if err != nil {
		if err.Error() == "todo no encontrado" {
			response.NotFound(w, err.Error())
			return
		}
		if err.Error() == "usuario no autorizado para marcar todos" {
			response.Forbidden(w, err.Error())
			return
		}
		response.InternalServerError(w, err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, completeResponse)
}

// DeleteTodo handles todo deletion
// @Summary Delete a todo
// @Description Delete a todo (only the creator can delete their todos)
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} models.SuccessResponse "Todo deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Todo not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}
	
	// Get todo ID from URL
	vars := mux.Vars(r)
	todoIDStr := vars["id"]
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "ID de todo inválido")
		return
	}
	
	// Delete todo
	err = h.todoService.DeleteTodo(todoID, userID)
	if err != nil {
		if err.Error() == "todo no encontrado" {
			response.NotFound(w, err.Error())
			return
		}
		if err.Error() == "solo el creador puede eliminar este todo" {
			response.Forbidden(w, err.Error())
			return
		}
		response.InternalServerError(w, err.Error())
		return
	}
	
	response.Success(w, "Todo eliminado exitosamente", nil)
}