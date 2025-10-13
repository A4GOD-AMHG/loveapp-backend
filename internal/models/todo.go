package models

import "time"

// Todo represents a todo item in the system
// @Description Todo item information
type Todo struct {
	ID              int64     `json:"id" db:"id" example:"1"`                                                 // Todo ID
	Title           string    `json:"title" db:"title" example:"Comprar comida"`                              // Todo title
	Description     string    `json:"description" db:"description" example:"Ir al supermercado y comprar..."` // Todo description
	CreatorID       int64     `json:"creator_id" db:"creator_id" example:"1"`                                 // Creator user ID
	CreatorUsername string    `json:"creator_username" db:"creator_username" example:"anyel"`                 // Creator username
	CompletedAnyel  bool      `json:"completed_anyel" db:"completed_anyel" example:"true"`                    // Completed by Anyel
	CompletedAlexis bool      `json:"completed_alexis" db:"completed_alexis" example:"false"`                 // Completed by Alexis
	IsCompleted     bool      `json:"is_completed" example:"false"`                                           // Overall completion status
	CreatedAt       time.Time `json:"created_at" db:"created_at" example:"2024-01-01T00:00:00Z"`              // Creation timestamp
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at" example:"2024-01-01T00:00:00Z"`              // Last update timestamp
}

// CreateTodoRequest represents the create todo request payload
// @Description Create todo request
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required" example:"Comprar comida"`                   // Todo title
	Description string `json:"description" example:"Ir al supermercado y comprar frutas y verduras"` // Todo description
}

// CreateTodoResponse represents the create todo response
// @Description Create todo response
type CreateTodoResponse struct {
	Message string `json:"message" example:"Todo creado exitosamente"` // Success message
	Todo    Todo   `json:"todo"`                                       // Created todo
}

// TodoListResponse represents the todo list response
// @Description Todo list response
type TodoListResponse struct {
	Message string `json:"message" example:"Todos obtenidos exitosamente"` // Success message
	Todos   []Todo `json:"todos"`                                          // List of todos
	Total   int    `json:"total" example:"5"`                              // Total count
}

// CompleteTodoResponse represents the complete todo response
// @Description Complete todo response
type CompleteTodoResponse struct {
	Message string `json:"message" example:"Todo marcado como completado"` // Success message
	Todo    Todo   `json:"todo"`                                           // Updated todo
}

// TodoStatus represents the possible todo statuses for filtering
type TodoStatus string

const (
	TodoStatusAll       TodoStatus = "all"       // All todos
	TodoStatusPending   TodoStatus = "pending"   // Pending todos
	TodoStatusCompleted TodoStatus = "completed" // Completed todos
)
