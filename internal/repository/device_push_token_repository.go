package repository

import (
	"database/sql"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

// DevicePushTokenRepository define las operaciones de persistencia para tokens push.
type DevicePushTokenRepository interface {
	Upsert(token *models.DevicePushToken) error
	DeleteByToken(userID int64, pushToken string) error
	FindByUserID(userID int64) ([]models.DevicePushToken, error)
}

type devicePushTokenRepository struct {
	db *sql.DB
}

// NewDevicePushTokenRepository crea un repositorio de tokens push.
func NewDevicePushTokenRepository() DevicePushTokenRepository {
	return &devicePushTokenRepository{db: database.DB}
}

func (r *devicePushTokenRepository) Upsert(token *models.DevicePushToken) error {
	if _, err := r.db.Exec(`DELETE FROM device_push_tokens WHERE push_token = ?`, token.PushToken); err != nil {
		return err
	}

	_, err := r.db.Exec(`
		INSERT INTO device_push_tokens (user_id, platform, push_token, device_name)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, platform, device_name)
		DO UPDATE SET push_token = excluded.push_token, updated_at = CURRENT_TIMESTAMP
	`, token.UserID, token.Platform, token.PushToken, token.DeviceName)

	return err
}

func (r *devicePushTokenRepository) DeleteByToken(userID int64, pushToken string) error {
	_, err := r.db.Exec(`DELETE FROM device_push_tokens WHERE user_id = ? AND push_token = ?`, userID, pushToken)
	return err
}

func (r *devicePushTokenRepository) FindByUserID(userID int64) ([]models.DevicePushToken, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, platform, push_token, device_name, created_at, updated_at
		FROM device_push_tokens
		WHERE user_id = ?
		ORDER BY id ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []models.DevicePushToken
	for rows.Next() {
		var token models.DevicePushToken
		if err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Platform,
			&token.PushToken,
			&token.DeviceName,
			&token.CreatedAt,
			&token.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}
