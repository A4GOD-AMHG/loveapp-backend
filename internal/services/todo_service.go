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

// CreateTodo creates a new todo
func (s *TodoService) CreateTodo(userID int64, username string, req *models.CreateTodoRequest) (*models.CreateTodoResponse, error) {
	// Check if user is allowed to create todos
	if !allowedCreators[username] {
		return nil, fmt.Errorf("usuario no autorizado para crear todos")
	}
	
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
		Message: "Todo creado exitosamente",
		Todo:    *createdTodo,
	}, nil
}

// ListTodos lists todos with optional filters
func (s *TodoService) ListTodos(statusStr, creatorIDStr string) (*models.TodoListResponse, error) {
	// Parse status
	var status models.TodoStatus
	switch statusStr {
	case "completed":
		status = models.TodoStatusCompleted
	case "pending":
		status = models.TodoStatusPending
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
	
	// Get todos
	todos, err := s.todoRepo.List(status, creatorID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener todos")
	}
	
	return &models.TodoListResponse{
		Message: "Todos obtenidos exitosamente",
		Todos:   todos,
		Total:   len(todos),
	}, nil
}

// CompleteTodo marks a todo as completed by a specific user
func (s *TodoService) CompleteTodo(todoID int64, username string) (*models.CompleteTodoResponse, error) {
	// Validate username
	if username != "anyel" && username != "alexis" {
		return nil, fmt.Errorf("usuario no autorizado para marcar todos")
	}
	
	// Check if todo exists
	_, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, fmt.Errorf("todo no encontrado")
	}
	
	// Update completion status
	updatedTodo, err := s.todoRepo.UpdateCompletion(todoID, username, true)
	if err != nil {
		return nil, fmt.Errorf("error al marcar todo como completado")
	}
	
	return &models.CompleteTodoResponse{
		Message: "Todo marcado como completado",
		Todo:    *updatedTodo,
	}, nil
}

// DeleteTodo deletes a todo (only by creator)
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