package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		db: database.DB,
	}
}

func (r *TodoRepository) Create(todo *models.Todo) (*models.Todo, error) {
	query := `
		INSERT INTO todos (title, description, creator_id, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	res, err := r.db.Exec(query, todo.Title, todo.Description, todo.CreatorID)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	todo.ID = id

	// Retrieve the created todo to get timestamps and other fields
	return r.FindByID(id)
}

func (r *TodoRepository) FindByID(id int64) (*models.Todo, error) {
	todo := &models.Todo{}
	query := `
		SELECT t.id, t.title, t.description, t.creator_id, u.username,
		       t.completed_anyel, t.completed_alexis, t.created_at, t.updated_at
		FROM todos t
		JOIN users u ON u.id = t.creator_id
		WHERE t.id = ?`

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

func (r *TodoRepository) GetTodos(status models.TodoStatus, creatorID *int64, requestingUsername, sortOrder, search string, limit, offset int) ([]models.Todo, int, error) {
	var args []interface{}

	baseQuery := `
		SELECT t.id, t.title, t.description, t.creator_id, u.username,
		       t.completed_anyel, t.completed_alexis, t.created_at, t.updated_at
		FROM todos t
		JOIN users u ON u.id = t.creator_id`

	conditions := []string{}

	if creatorID != nil {
		conditions = append(conditions, "t.creator_id = ?")
		args = append(args, *creatorID)
	}

	if search != "" {
		conditions = append(conditions, "(t.title LIKE ? OR t.description LIKE ?)")
		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam)
	}

	switch status {
	case models.TodoStatusCompleted:
		conditions = append(conditions, "t.completed_anyel = 1 AND t.completed_alexis = 1")
	case models.TodoStatusIncompleted:
		conditions = append(conditions, "(t.completed_anyel = 0 OR t.completed_alexis = 0)")
	case models.TodoStatusCompletedByMe:
		if requestingUsername == "anyel" {
			conditions = append(conditions, "t.completed_anyel = 1")
		} else if requestingUsername == "alexis" {
			conditions = append(conditions, "t.completed_alexis = 1")
		}
	}

	countQuery := `SELECT COUNT(*) FROM todos t JOIN users u ON u.id = t.creator_id`
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if strings.ToLower(sortOrder) == "asc" {
		query += " ORDER BY t.created_at ASC"
	} else {
		query += " ORDER BY t.created_at DESC"
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	todos := []models.Todo{}
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
			fmt.Printf("Error scanning todo row: %v\n", err)
			continue
		}
		todo.IsCompleted = todo.CompletedAnyel && todo.CompletedAlexis
		todos = append(todos, todo)
	}

	return todos, total, nil
}

func (r *TodoRepository) Update(todo *models.Todo) (*models.Todo, error) {
	query := `UPDATE todos SET title = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, todo.Title, todo.Description, todo.ID)
	if err != nil {
		return nil, err
	}
	return r.FindByID(todo.ID)
}

func (r *TodoRepository) UpdateCompletion(todoID int64, username string, completed bool) (*models.Todo, error) {
	var query string

	switch username {
	case "anyel":
		query = `UPDATE todos SET completed_anyel = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	case "alexis":
		query = `UPDATE todos SET completed_alexis = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
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

	return r.FindByID(todoID)
}

func (r *TodoRepository) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = ?`

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

func (r *TodoRepository) GetCreatorID(todoID int64) (int64, error) {
	var creatorID int64
	query := `SELECT creator_id FROM todos WHERE id = ?`

	err := r.db.QueryRow(query, todoID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("todo not found")
		}
		return 0, err
	}

	return creatorID, nil
}
