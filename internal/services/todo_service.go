// Paquete services implementa la lógica de negocio de la aplicación.
package services

import (
	"fmt"
	"strconv"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
)

// TodoService encapsula la lógica de negocio relacionada con la gestión de tareas.
type TodoService struct {
	todoRepo *repository.TodoRepository // Repositorio de tareas para acceso a datos
	userRepo *repository.UserRepository // Repositorio de usuarios (para validaciones futuras)
}

// NewTodoService crea y retorna una nueva instancia de TodoService
// con sus repositorios de tareas y usuarios ya inicializados.
func NewTodoService() *TodoService {
	return &TodoService{
		todoRepo: repository.NewTodoRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// allowedCreators define qué usuarios del sistema pueden crear tareas.
// En esta aplicación solo anyel y alexis tienen acceso.
var allowedCreators = map[string]bool{
	"anyel":  true,
	"alexis": true,
}

// CreateTodo crea una nueva tarea en el sistema para el usuario especificado.
// Valida que el título no esté vacío antes de persistir.
// Retorna la respuesta con el mensaje de éxito y los datos de la tarea creada.
func (s *TodoService) CreateTodo(userID int64, username string, req *models.CreateTodoRequest) (*models.CreateTodoResponse, error) {
	// Validar que el título de la tarea no esté vacío
	if req.Title == "" {
		return nil, fmt.Errorf("el título es requerido")
	}

	// Construir la entidad Todo con los datos del usuario creador
	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   userID,
	}

	// Persistir la tarea en la base de datos
	createdTodo, err := s.todoRepo.Create(todo)
	if err != nil {
		return nil, fmt.Errorf("error al crear todo")
	}

	return &models.CreateTodoResponse{
		Message: "¡Tarea creada con éxito! 🚀",
		Todo:    *createdTodo,
	}, nil
}

// GetTodos retorna una lista paginada de tareas con soporte para múltiples filtros.
// Parámetros recibidos como strings (provenientes de query params) y convertidos internamente.
// Soporta filtros por: estado, creador, búsqueda de texto y ordenamiento.
func (s *TodoService) GetTodos(statusStr, creatorIDStr, username, sortOrder, search, pageStr, limitStr string) (*models.TodoListResponse, error) {
	// Convertir el string de estado al tipo TodoStatus correspondiente
	var status models.TodoStatus
	switch statusStr {
	case "completed":
		status = models.TodoStatusCompleted
	case "incompleted":
		status = models.TodoStatusIncompleted
	case "completed_by_me":
		status = models.TodoStatusCompletedByMe
	default:
		// Si no se especifica o es inválido, retornar todas las tareas
		status = models.TodoStatusAll
	}

	// Convertir y validar el filtro de ID de creador (opcional)
	var creatorID *int64
	if creatorIDStr != "" {
		id, err := strconv.ParseInt(creatorIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ID de creador inválido")
		}
		creatorID = &id
	}

	// Parsear y validar los parámetros de paginación con valores predeterminados
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Consultar las tareas al repositorio con todos los filtros aplicados
	todos, total, err := s.todoRepo.GetTodos(status, creatorID, username, sortOrder, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al obtener todos")
	}

	// Calcular la última página según el total de resultados
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

// UpdateTodo actualiza el título y la descripción de una tarea existente.
// Solo el creador de la tarea puede modificarla.
// Retorna error si la tarea no existe o si el usuario no es el creador.
func (s *TodoService) UpdateTodo(todoID int64, userID int64, req *models.UpdateTodoRequest) (*models.Todo, error) {
	// Verificar que la tarea existe en la base de datos
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, fmt.Errorf("todo no encontrado")
	}

	// Verificar que el usuario solicitante es el creador de la tarea
	if todo.CreatorID != userID {
		return nil, fmt.Errorf("solo el creador puede editar este todo")
	}

	// Aplicar los cambios al modelo antes de persistir
	todo.Title = req.Title
	todo.Description = req.Description

	return s.todoRepo.Update(todo)
}

// UpdateTodoStatus actualiza el estado de completado de una tarea para un usuario específico.
// Solo "anyel" y "alexis" pueden marcar tareas.
// Una vez que ambos han completado la tarea (IsCompleted = true), no se puede desmarcar.
func (s *TodoService) UpdateTodoStatus(todoID int64, username string, completed bool) (*models.CompleteTodoResponse, error) {
	// Validar que el usuario sea uno de los dos usuarios permitidos del sistema
	if username != "anyel" && username != "alexis" {
		return nil, fmt.Errorf("usuario no autorizado para marcar todos")
	}

	// Verificar que la tarea existe
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, fmt.Errorf("todo no encontrado")
	}

	// Regla de negocio: una vez completada por ambos, la tarea no puede desmarcarse
	if todo.IsCompleted {
		if !completed {
			return nil, fmt.Errorf("la tarea ya está completada por ambos y no se puede desmarcar")
		}
		// Si intenta marcarla como completada nuevamente cuando ya lo está, no hay error
	}

	// Actualizar el estado de completado en la base de datos para este usuario
	updatedTodo, err := s.todoRepo.UpdateCompletion(todoID, username, completed)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estado del todo")
	}

	// Personalizar el mensaje según si ambos completaron la tarea o solo uno
	message := "Estado actualizado correctamente"
	if updatedTodo.IsCompleted {
		message = "¡Bien hecho! Tarea completada por ambos. 🎉"
	}

	return &models.CompleteTodoResponse{
		Message: message,
		Todo:    *updatedTodo,
	}, nil
}

// DeleteTodo elimina permanentemente una tarea del sistema.
// Solo el creador de la tarea puede eliminarla.
// Retorna error si la tarea no existe o si el usuario no es el creador.
func (s *TodoService) DeleteTodo(todoID int64, userID int64) error {
	// Obtener el ID del creador de la tarea para verificar permisos
	creatorID, err := s.todoRepo.GetCreatorID(todoID)
	if err != nil {
		return fmt.Errorf("todo no encontrado")
	}

	// Verificar que el usuario solicitante sea el creador
	if creatorID != userID {
		return fmt.Errorf("solo el creador puede eliminar este todo")
	}

	// Eliminar la tarea de la base de datos
	err = s.todoRepo.Delete(todoID)
	if err != nil {
		return fmt.Errorf("error al eliminar todo")
	}

	return nil
}
