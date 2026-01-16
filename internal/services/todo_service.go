package services

import (
	"fmt"
	"strconv"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
)

// TodoService handles todo business logic
type TodoService struct {
	todoRepo *repository.TodoRepository
	userRepo *repository.UserRepository
}

// NewTodoService creates a new todo service
func NewTodoService() *TodoService {
	return &TodoService{
		todoRepo: repository.NewTodoRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// allowedCreators defines which users can create todos
var allowedCreators = map[string]bool{
	"anyel":  true,
	"alexis": true,
}

// CreateTodo crea una nueva tarea
func (s *TodoService) CreateTodo(userID int64, username string, req *models.CreateTodoRequest) (*models.CreateTodoResponse, error) {
	// Eliminar restricción explícita de usuarios, asumir que el middleware de autenticación maneja el acceso.
	// O mantener si solo estos usuarios deben crear, pero el usuario dijo que "ambos usuarios pueden crear".
	
	// Validate input
	if req.Title == "" {
		return nil, fmt.Errorf("el título es requerido")
	}
	
	// Create todo
	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   userID,
	}
	
	createdTodo, err := s.todoRepo.Create(todo)
	if err != nil {
		return nil, fmt.Errorf("error al crear todo")
	}
	
	return &models.CreateTodoResponse{
		Message: "¡Tarea creada con éxito! 🚀",
		Todo:    *createdTodo,
	}, nil
}

// GetTodos lista las tareas con filtros y paginación
func (s *TodoService) GetTodos(statusStr, creatorIDStr, username, sortOrder, search, pageStr, limitStr string) (*models.TodoListResponse, error) {
	// Parse status
	var status models.TodoStatus
	switch statusStr {
	case "completed":
		status = models.TodoStatusCompleted
	case "incompleted":
		status = models.TodoStatusIncompleted
	case "completed_by_me":
		status = models.TodoStatusCompletedByMe
	default:
		status = models.TodoStatusAll
	}
	
	// Parse creator ID
	var creatorID *int64
	if creatorIDStr != "" {
		id, err := strconv.ParseInt(creatorIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ID de creador inválido")
		}
		creatorID = &id
	}

	// Parse pagination
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	
	// Get todos
	todos, total, err := s.todoRepo.GetTodos(status, creatorID, username, sortOrder, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al obtener todos")
	}

	lastPage := total / limit
	if total%limit != 0 {
		lastPage++
	}

	return &models.TodoListResponse{
		Message:  "Aquí están todas las tareas. ¡Vamos a completarlas! 💪",
		Todos:    todos,
		Total:    total,
		Page:     page,
		PerPage:  limit,
		LastPage: lastPage,
	}, nil
}

// UpdateTodo actualiza una tarea (solo título y descripción, solo por creador)
func (s *TodoService) UpdateTodo(todoID int64, userID int64, req *models.UpdateTodoRequest) (*models.Todo, error) {
	// Check if todo exists
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, fmt.Errorf("todo no encontrado")
	}
	
	// Check if user is the creator
	if todo.CreatorID != userID {
		return nil, fmt.Errorf("solo el creador puede editar este todo")
	}
	
	todo.Title = req.Title
	todo.Description = req.Description
	
	return s.todoRepo.Update(todo)
}

// UpdateTodoStatus actualiza el estado de la tarea para un usuario específico
func (s *TodoService) UpdateTodoStatus(todoID int64, username string, completed bool) (*models.CompleteTodoResponse, error) {
	// Validate username (hardcoded check as per previous logic, ensuring only valid users)
	if username != "anyel" && username != "alexis" {
		return nil, fmt.Errorf("usuario no autorizado para marcar todos")
	}
	
	// Check if todo exists
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, fmt.Errorf("todo no encontrado")
	}

	// Logic: If IsCompleted is already true (both completed), prevent unchecking.
	// "ese iscompleted una vez marcado no se puede desmarcar"
	if todo.IsCompleted {
		// If both have completed it, we shouldn't allow changing the state back to false
		// for either user.
		// However, we need to check if the incoming request is trying to set completed to false.
		if !completed {
			// User is trying to uncheck
			return nil, fmt.Errorf("la tarea ya está completada por ambos y no se puede desmarcar")
		}
		// If user sets completed to true again, it's fine, it stays true.
	}
	
	// Update completion status in DB
	updatedTodo, err := s.todoRepo.UpdateCompletion(todoID, username, completed)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estado del todo")
	}

	// Logic: Task is completed globally only if both have marked it (handled in model/repo usually or here)
	// The repo returns the updated todo with IsCompleted calculated.
	// If the user wants a message "Todo completed" only when fully completed, we can check updatedTodo.IsCompleted.
	
	message := "Estado actualizado correctamente"
	if updatedTodo.IsCompleted {
		message = "¡Bien hecho! Tarea completada por ambos. 🎉"
	}
	
	return &models.CompleteTodoResponse{
		Message: message,
		Todo:    *updatedTodo,
	}, nil
}

// DeleteTodo elimina una tarea (solo por el creador)
func (s *TodoService) DeleteTodo(todoID int64, userID int64) error {
	// Check if todo exists and get creator ID
	creatorID, err := s.todoRepo.GetCreatorID(todoID)
	if err != nil {
		return fmt.Errorf("todo no encontrado")
	}
	
	// Check if user is the creator
	if creatorID != userID {
		return fmt.Errorf("solo el creador puede eliminar este todo")
	}
	
	// Delete todo
	err = s.todoRepo.Delete(todoID)
	if err != nil {
		return fmt.Errorf("error al eliminar todo")
	}
	
	return nil
}