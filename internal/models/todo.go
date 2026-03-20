// Paquete models define las estructuras de datos utilizadas en toda la aplicación.
package models

import "time"

// Todo representa un elemento de tarea compartida entre los dos usuarios del sistema.
// Cada tarea tiene un estado de completado independiente por usuario (anyel y alexis).
// @Description Información del elemento de tarea
type Todo struct {
	ID              int64     `json:"id" db:"id" example:"1"`                                                  // Identificador único de la tarea
	Title           string    `json:"title" db:"title" example:"Comprar comida"`                               // Título descriptivo de la tarea
	Description     string    `json:"description" db:"description" example:"Ir al supermercado y comprar..."` // Descripción detallada de la tarea
	CreatorID       int64     `json:"creator_id" db:"creator_id" example:"1"`                                  // ID del usuario que creó la tarea
	CreatorUsername string    `json:"creator_username" db:"creator_username" example:"anyel"`                  // Nombre de usuario del creador
	CompletedAnyel  bool      `json:"completed_anyel" db:"completed_anyel" example:"true"`                     // Indica si Anyel marcó la tarea como completada
	CompletedAlexis bool      `json:"completed_alexis" db:"completed_alexis" example:"false"`                  // Indica si Alexis marcó la tarea como completada
	IsCompleted     bool      `json:"is_completed" example:"false"`                                            // true solo cuando ambos usuarios han completado la tarea
	CreatedAt       time.Time `json:"created_at" db:"created_at" example:"2024-01-01T00:00:00Z"`               // Fecha y hora de creación
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at" example:"2024-01-01T00:00:00Z"`               // Fecha y hora de la última actualización
}

// CreateTodoRequest representa los datos necesarios para crear una nueva tarea.
// @Description Solicitud de creación de tarea
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required" example:"Comprar comida"`                               // Título de la tarea (obligatorio)
	Description string `json:"description" example:"Ir al supermercado y comprar frutas y verduras"` // Descripción opcional de la tarea
}

// CreateTodoResponse representa la respuesta retornada al crear una tarea exitosamente.
// @Description Respuesta de creación de tarea
type CreateTodoResponse struct {
	Message string `json:"message" example:"¡Tarea creada con éxito! 🚀"` // Mensaje de confirmación
	Todo    Todo   `json:"todo"`                                          // Datos completos de la tarea recién creada
}

// TodoListResponse representa la respuesta paginada al listar tareas.
// @Description Respuesta de lista de tareas
type TodoListResponse struct {
	Message  string `json:"message" example:"Aquí están todas las tareas. ¡Vamos a completarlas! 💪"` // Mensaje informativo
	Todos    []Todo `json:"todos"`                                                                     // Lista de tareas encontradas
	Total    int    `json:"total" example:"5"`                                                         // Total de tareas que coinciden con el filtro
	Page     int    `json:"page" example:"1"`                                                          // Página actual
	PerPage  int    `json:"per_page" example:"10"`                                                     // Cantidad de tareas por página
	LastPage int    `json:"last_page" example:"1"`                                                     // Última página disponible
}

// CompleteTodoResponse representa la respuesta al actualizar el estado de una tarea.
// @Description Respuesta de finalización de tarea
type CompleteTodoResponse struct {
	Message string `json:"message" example:"¡Bien hecho! Tarea marcada como completada. 🎉"` // Mensaje de resultado
	Todo    Todo   `json:"todo"`                                                             // Estado actualizado de la tarea
}

// UpdateTodoRequest representa los datos para actualizar el título y descripción de una tarea.
// @Description Solicitud de actualización de tarea
type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required" example:"Comprar comida"`                               // Nuevo título de la tarea (obligatorio)
	Description string `json:"description" example:"Ir al supermercado y comprar frutas y verduras"` // Nueva descripción de la tarea
}

// UpdateTodoStatusRequest representa la solicitud para marcar o desmarcar una tarea como completada.
// @Description Solicitud de actualización de estado de tarea
type UpdateTodoStatusRequest struct {
	Completed bool `json:"completed" example:"true"` // true para marcar como completada, false para desmarcar
}

// TodoStatus define los posibles valores de filtro de estado para listar tareas.
type TodoStatus string

const (
	// TodoStatusAll retorna todas las tareas sin filtro de estado
	TodoStatusAll TodoStatus = "all"
	// TodoStatusCompleted retorna solo las tareas completadas por ambos usuarios
	TodoStatusCompleted TodoStatus = "completed"
	// TodoStatusIncompleted retorna las tareas pendientes de al menos un usuario
	TodoStatusIncompleted TodoStatus = "incompleted"
	// TodoStatusCompletedByMe retorna las tareas completadas por el usuario actual
	TodoStatusCompletedByMe TodoStatus = "completed_by_me"
)
