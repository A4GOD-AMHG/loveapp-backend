// Paquete repository implementa el acceso a datos de la aplicación.
package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

// TodoRepository gestiona las operaciones de base de datos para las tareas (todos).
type TodoRepository struct {
	db *sql.DB // Conexión a la base de datos SQLite
}

// NewTodoRepository crea y retorna una nueva instancia de TodoRepository
// conectada a la base de datos global de la aplicación.
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		db: database.DB,
	}
}

// Create inserta una nueva tarea en la base de datos y retorna el registro completo
// con los campos generados automáticamente (ID, timestamps, username del creador).
func (r *TodoRepository) Create(todo *models.Todo) (*models.Todo, error) {
	query := `
		INSERT INTO todos (title, description, creator_id, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	res, err := r.db.Exec(query, todo.Title, todo.Description, todo.CreatorID)
	if err != nil {
		return nil, err
	}

	// Obtener el ID generado por la base de datos
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	todo.ID = id

	// Recargar el registro completo desde la base de datos para obtener timestamps y datos JOIN
	return r.FindByID(id)
}

// FindByID busca una tarea por su ID y hace JOIN con la tabla de usuarios
// para obtener el nombre de usuario del creador.
// Retorna error si la tarea no existe.
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

	// Calcular el campo virtual IsCompleted: true solo si ambos usuarios completaron la tarea
	todo.IsCompleted = todo.CompletedAnyel && todo.CompletedAlexis
	return todo, nil
}

// GetTodos retorna una lista paginada de tareas con soporte para filtros múltiples.
// Parámetros:
//   - status: filtra por estado de completado (all, completed, incompleted, completed_by_me)
//   - creatorID: filtra por ID del creador (nil = sin filtro)
//   - requestingUsername: nombre del usuario que realiza la consulta (para filtro completed_by_me)
//   - sortOrder: orden de resultados por fecha ("asc" o "desc")
//   - search: búsqueda de texto en título y descripción
//   - limit, offset: parámetros de paginación
//
// Retorna la lista de tareas, el total de resultados y un error si ocurre.
func (r *TodoRepository) GetTodos(status models.TodoStatus, creatorID *int64, requestingUsername, sortOrder, search string, limit, offset int) ([]models.Todo, int, error) {
	var args []interface{}

	// Consulta base con JOIN para obtener el username del creador
	baseQuery := `
		SELECT t.id, t.title, t.description, t.creator_id, u.username,
		       t.completed_anyel, t.completed_alexis, t.created_at, t.updated_at
		FROM todos t
		JOIN users u ON u.id = t.creator_id`

	conditions := []string{}

	// Filtro por creador específico
	if creatorID != nil {
		conditions = append(conditions, "t.creator_id = ?")
		args = append(args, *creatorID)
	}

	// Filtro de búsqueda de texto en título o descripción
	if search != "" {
		conditions = append(conditions, "(t.title LIKE ? OR t.description LIKE ?)")
		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam)
	}

	// Filtro por estado de completado
	switch status {
	case models.TodoStatusCompleted:
		// Solo tareas completadas por ambos usuarios
		conditions = append(conditions, "t.completed_anyel = 1 AND t.completed_alexis = 1")
	case models.TodoStatusIncompleted:
		// Tareas donde al menos un usuario no ha completado
		conditions = append(conditions, "(t.completed_anyel = 0 OR t.completed_alexis = 0)")
	case models.TodoStatusCompletedByMe:
		// Tareas completadas específicamente por el usuario que consulta
		if requestingUsername == "anyel" {
			conditions = append(conditions, "t.completed_anyel = 1")
		} else if requestingUsername == "alexis" {
			conditions = append(conditions, "t.completed_alexis = 1")
		}
	}

	// Construir y ejecutar la consulta de conteo total (sin paginación)
	countQuery := `SELECT COUNT(*) FROM todos t JOIN users u ON u.id = t.creator_id`
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Construir la consulta principal con los mismos filtros
	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Aplicar ordenamiento por fecha de creación
	if strings.ToLower(sortOrder) == "asc" {
		query += " ORDER BY t.created_at ASC"
	} else {
		query += " ORDER BY t.created_at DESC"
	}

	// Aplicar paginación
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Mapear cada fila al modelo Todo
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
		// Calcular el campo virtual IsCompleted
		todo.IsCompleted = todo.CompletedAnyel && todo.CompletedAlexis
		todos = append(todos, todo)
	}

	return todos, total, nil
}

// Update actualiza el título y la descripción de una tarea existente.
// Retorna el registro actualizado con los nuevos timestamps.
func (r *TodoRepository) Update(todo *models.Todo) (*models.Todo, error) {
	query := `UPDATE todos SET title = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, todo.Title, todo.Description, todo.ID)
	if err != nil {
		return nil, err
	}
	// Recargar desde la base de datos para retornar el estado actualizado
	return r.FindByID(todo.ID)
}

// UpdateCompletion actualiza el estado de completado de una tarea para un usuario específico.
// Solo acepta "anyel" o "alexis" como valores de username.
// Retorna error si la tarea no existe o el username es inválido.
func (r *TodoRepository) UpdateCompletion(todoID int64, username string, completed bool) (*models.Todo, error) {
	var query string

	// Seleccionar la columna de completado según el usuario
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

	// Verificar que la tarea existiera en la base de datos
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("todo not found")
	}

	// Retornar el estado actualizado de la tarea
	return r.FindByID(todoID)
}

// Delete elimina permanentemente una tarea de la base de datos por su ID.
// Retorna error si la tarea no existe.
func (r *TodoRepository) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	// Verificar que la tarea existiera para informar al llamador
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

// GetCreatorID retorna el ID del usuario creador de una tarea.
// Útil para verificar permisos de modificación/eliminación sin cargar el registro completo.
// Retorna error si la tarea no existe.
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
