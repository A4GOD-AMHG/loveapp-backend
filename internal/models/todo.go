package models

import "time"

// Todo representa un elemento de tarea en el sistema
// @Description Información del elemento de tarea
type Todo struct {
	ID              int64     `json:"id" db:"id" example:"1"`
	Title           string    `json:"title" db:"title" example:"Comprar comida"`
	Description     string    `json:"description" db:"description" example:"Ir al supermercado y comprar..."`
	CreatorID       int64     `json:"creator_id" db:"creator_id" example:"1"`
	CreatorUsername string    `json:"creator_username" db:"creator_username" example:"anyel"`
	CompletedAnyel  bool      `json:"completed_anyel" db:"completed_anyel" example:"true"`
	CompletedAlexis bool      `json:"completed_alexis" db:"completed_alexis" example:"false"`
	IsCompleted     bool      `json:"is_completed" example:"false"`
	CreatedAt       time.Time `json:"created_at" db:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// CreateTodoRequest representa la carga útil de solicitud de creación de tarea
// @Description Solicitud de creación de tarea
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required" example:"Comprar comida"`
	Description string `json:"description" example:"Ir al supermercado y comprar frutas y verduras"`
}

// CreateTodoResponse representa la respuesta de creación de tarea
// @Description Respuesta de creación de tarea
type CreateTodoResponse struct {
	Message string `json:"message" example:"¡Tarea creada con éxito! 🚀"`
	Todo    Todo   `json:"todo"`
}

// TodoListResponse representa la respuesta de lista de tareas
// @Description Respuesta de lista de tareas
type TodoListResponse struct {
	Message  string `json:"message" example:"Aquí están todas las tareas. ¡Vamos a completarlas! 💪"`
	Todos    []Todo `json:"todos"`
	Total    int    `json:"total" example:"5"`
	Page     int    `json:"page" example:"1"`
	PerPage  int    `json:"per_page" example:"10"`
	LastPage int    `json:"last_page" example:"1"`
}

// CompleteTodoResponse representa la respuesta de finalización de tarea
// @Description Respuesta de finalización de tarea
type CompleteTodoResponse struct {
	Message string `json:"message" example:"¡Bien hecho! Tarea marcada como completada. 🎉"`
	Todo    Todo   `json:"todo"`
}

// UpdateTodoRequest representa la carga útil de solicitud de actualización de tarea
// @Description Solicitud de actualización de tarea
type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required" example:"Comprar comida"`
	Description string `json:"description" example:"Ir al supermercado y comprar frutas y verduras"`
}

// UpdateTodoStatusRequest representa la carga útil de solicitud de actualización de estado de tarea
// @Description Solicitud de actualización de estado de tarea
type UpdateTodoStatusRequest struct {
	Completed bool `json:"completed" example:"true"`
}

// TodoStatus representa los posibles estados de tarea para filtrar
type TodoStatus string

const (
	TodoStatusAll           TodoStatus = "all"
	TodoStatusCompleted     TodoStatus = "completed"
	TodoStatusIncompleted   TodoStatus = "incompleted"
	TodoStatusCompletedByMe TodoStatus = "completed_by_me"
)
