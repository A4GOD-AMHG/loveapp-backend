package repository

import (
	"database/sql"
	"fmt"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

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

func (r *UserRepository) UpdatePassword(userID int64, hashedPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	result, err := r.db.Exec(query, hashedPassword, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

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
