// Paquete repository implementa el acceso a datos de la aplicación.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

// UserRepository gestiona las operaciones de base de datos para los usuarios.
type UserRepository struct {
	db *sql.DB // Conexión a la base de datos SQLite
}

// NewUserRepository crea y retorna una nueva instancia de UserRepository
// conectada a la base de datos global de la aplicación.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

// FindByUsername busca un usuario por su nombre de usuario, incluyendo el hash de la contraseña.
// Se usa principalmente en el proceso de autenticación (login).
// Retorna error si el usuario no existe.
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, username, password, created_at, updated_at FROM users WHERE username = ?`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

// FindByID busca un usuario por su ID sin incluir la contraseña.
// Se usa para obtener datos del usuario autenticado desde el contexto o el middleware.
// Retorna error si el usuario no existe.
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, username, created_at, updated_at FROM users WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UpdatePassword actualiza el hash de la contraseña de un usuario.
// También actualiza automáticamente el campo updated_at (mediante el trigger de base de datos).
// Retorna error si el usuario no existe o si la actualización falla.
func (r *UserRepository) UpdatePassword(userID int64, hashedPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	result, err := r.db.Exec(query, hashedPassword, userID)
	if err != nil {
		return err
	}

	// Verificar que el usuario existiera para detectar IDs inválidos
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetPasswordHash retorna el hash de la contraseña almacenado para un usuario dado su ID.
// Útil para verificar la contraseña actual antes de permitir cambios.
// Retorna error si el usuario no existe.
func (r *UserRepository) GetPasswordHash(userID int64) (string, error) {
	var passwordHash string
	query := `SELECT password FROM users WHERE id = ?`

	err := r.db.QueryRow(query, userID).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}

	return passwordHash, nil
}

// GetOtherUser retorna el único usuario del sistema que NO sea el usuario actual.
// Dado que la aplicación está diseñada para exactamente dos usuarios, este método
// determina el destinatario de los mensajes automáticamente sin necesidad de especificarlo.
// Retorna error si no existe otro usuario registrado.
func (r *UserRepository) GetOtherUser(currentUserID uint) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, username, created_at, updated_at FROM users WHERE id != ? LIMIT 1`

	err := r.db.QueryRow(query, currentUserID).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("other user not found")
		}
		return nil, err
	}

	return user, nil
}
