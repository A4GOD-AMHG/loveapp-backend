// Paquete services implementa la lógica de negocio de la aplicación.
package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
)

// TodoService encapsula la lógica de negocio relacionada con la gestión de tareas.
type TodoService struct {
	todoRepo *repository.TodoRepository
}

// NewTodoService crea una instancia de TodoService.
func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

// CreateTodo crea una nueva tarea para el usuario especificado.
func (s *TodoService) CreateTodo(userID int64, req *models.CreateTodoRequest) (*models.CreateTodoResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrInvalidInput
	}

	createdTodo, err := s.todoRepo.Create(&models.Todo{
		Title:       title,
		Description: req.Description,
		CreatorID:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear todo: %w", err)
	}

	return &models.CreateTodoResponse{
		Message: "¡Tarea creada con éxito! 🚀",
		Todo:    *createdTodo,
	}, nil
}

// GetTodos retorna una lista paginada de tareas. Los parámetros llegan como
// strings desde la query y se convierten aquí con valores predeterminados.
func (s *TodoService) GetTodos(statusStr, creatorIDStr, username, sortOrder, search, pageStr, limitStr string) (*models.TodoListResponse, error) {
	status := models.TodoStatus(statusStr)
	switch status {
	case models.TodoStatusCompleted, models.TodoStatusIncompleted, models.TodoStatusCompletedByMe:
	default:
		status = models.TodoStatusAll
	}

	var creatorID *int64
	if creatorIDStr != "" {
		id, err := strconv.ParseInt(creatorIDStr, 10, 64)
		if err != nil {
			return nil, ErrInvalidInput
		}
		creatorID = &id
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	todos, total, err := s.todoRepo.GetTodos(status, creatorID, username, sortOrder, search, limit, (page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("error al obtener todos: %w", err)
	}

	return &models.TodoListResponse{
		Message:  "Aquí están todas las tareas. ¡Vamos a completarlas! 💪",
		Todos:    todos,
		Total:    total,
		Page:     page,
		PerPage:  limit,
		LastPage: (total + limit - 1) / limit,
	}, nil
}

// UpdateTodo actualiza título y descripción; solo el creador puede hacerlo.
func (s *TodoService) UpdateTodo(todoID int64, userID int64, req *models.UpdateTodoRequest) (*models.Todo, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrInvalidInput
	}

	todo, err := s.find(todoID)
	if err != nil {
		return nil, err
	}
	if todo.CreatorID != userID {
		return nil, ErrForbidden
	}

	todo.Title = title
	todo.Description = req.Description
	return s.todoRepo.Update(todo)
}

// UpdateTodoStatus marca o desmarca la tarea para el usuario.
// Una vez completada por ambos, no se puede desmarcar.
func (s *TodoService) UpdateTodoStatus(todoID int64, username string, completed bool) (*models.CompleteTodoResponse, error) {
	if username != "anyel" && username != "alexis" {
		return nil, ErrForbidden
	}

	todo, err := s.find(todoID)
	if err != nil {
		return nil, err
	}
	if todo.IsCompleted && !completed {
		return nil, ErrAlreadyCompleted
	}

	updatedTodo, err := s.todoRepo.UpdateCompletion(todoID, username, completed)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estado del todo: %w", err)
	}

	message := "Estado actualizado correctamente"
	if updatedTodo.IsCompleted {
		message = "¡Bien hecho! Tarea completada por ambos. 🎉"
	}

	return &models.CompleteTodoResponse{
		Message: message,
		Todo:    *updatedTodo,
	}, nil
}

// DeleteTodo elimina una tarea; solo el creador puede hacerlo.
func (s *TodoService) DeleteTodo(todoID int64, userID int64) error {
	creatorID, err := s.todoRepo.GetCreatorID(todoID)
	if err != nil {
		return notFoundOr(err)
	}
	if creatorID != userID {
		return ErrForbidden
	}
	if err := s.todoRepo.Delete(todoID); err != nil {
		return notFoundOr(err)
	}
	return nil
}

func (s *TodoService) find(todoID int64) (*models.Todo, error) {
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, notFoundOr(err)
	}
	return todo, nil
}

// notFoundOr traduce el "no encontrado" del repositorio al error de dominio.
func notFoundOr(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
