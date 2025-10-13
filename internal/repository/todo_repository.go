package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

// TodoRepository handles todo data operations
type TodoRepository struct {
	db *sql.DB
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		db: database.DB,
	}
}

// Create creates a new todo
func (r *TodoRepository) Create(todo *models.Todo) (*models.Todo, error) {
	query := `
		INSERT INTO todos (title, description, creator_id) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`
	
	err := r.db.QueryRow(query, todo.Title, todo.Description, todo.CreatorID).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Get the creator username
	userRepo := NewUserRepository()
	user, err := userRepo.FindByID(todo.CreatorID)
	if err != nil {
		return nil, err
	}
	todo.CreatorUsername = user.Username
	
	return todo, nil
}

// FindByID finds a todo by ID
func (r *TodoRepository) FindByID(id int64) (*models.Todo, error) {
	todo := &models.Todo{}
	query := `
		SELECT t.id, t.title, t.description, t.creator_id, u.username, 
		       t.completed_anyel, t.completed_alexis, t.created_at, t.updated_at
		FROM todos t 
		JOIN users u ON u.id = t.creator_id 
		WHERE t.id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.CreatorID,
		&todo.CreatorUsername,
		&todo.CompletedAnyel,
		&todo.CompletedAlexis,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, err
	}
	
	todo.IsCompleted = todo.CompletedAnyel && todo.CompletedAlexis
	return todo, nil
}

// List lists todos with optional filters
func (r *TodoRepository) List(status models.TodoStatus, creatorID *int64) ([]models.Todo, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	baseQuery := `
		SELECT t.id, t.title, t.description, t.creator_id, u.username, 
		       t.completed_anyel, t.completed_alexis, t.created_at, t.updated_at
		FROM todos t 
		JOIN users u ON u.id = t.creator_id`
	
	// Add creator filter
	if creatorID != nil {
		conditions = append(conditions, fmt.Sprintf("t.creator_id = $%d", argIndex))
		args = append(args, *creatorID)
		argIndex++
	}
	
	// Build the query with conditions
	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY t.created_at DESC"
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.CreatorID,
			&todo.CreatorUsername,
			&todo.CompletedAnyel,
			&todo.CompletedAlexis,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			continue
		}
		
		todo.IsCompleted = todo.CompletedAnyel && todo.CompletedAlexis
		
		// Apply status filter
		switch status {
		case models.TodoStatusCompleted:
			if !todo.IsCompleted {
				continue
			}
		case models.TodoStatusPending:
			if todo.IsCompleted {
				continue
			}
		case models.TodoStatusAll:
			// Include all todos
		}
		
		todos = append(todos, todo)
	}
	
	return todos, nil
}

// UpdateCompletion updates the completion status for a specific user
func (r *TodoRepository) UpdateCompletion(todoID int64, username string, completed bool) (*models.Todo, error) {
	var query string
	
	switch username {
	case "anyel":
		query = `UPDATE todos SET completed_anyel = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	case "alexis":
		query = `UPDATE todos SET completed_alexis = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	default:
		return nil, fmt.Errorf("invalid username for completion update")
	}
	
	result, err := r.db.Exec(query, completed, todoID)
	if err != nil {
		return nil, err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	
	if rowsAffected == 0 {
		return nil, fmt.Errorf("todo not found")
	}
	
	// Return the updated todo
	return r.FindByID(todoID)
}

// Delete deletes a todo
func (r *TodoRepository) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}
	
	return nil
}

// GetCreatorID gets the creator ID of a todo
func (r *TodoRepository) GetCreatorID(todoID int64) (int64, error) {
	var creatorID int64
	query := `SELECT creator_id FROM todos WHERE id = $1`
	
	err := r.db.QueryRow(query, todoID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("todo not found")
		}
		return 0, err
	}
	
	return creatorID, nil
}