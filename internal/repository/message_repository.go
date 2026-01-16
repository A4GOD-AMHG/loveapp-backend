package repository

import (
	"database/sql"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

type MessageRepository interface {
	Create(message *models.Message) (int64, error)
	FindByID(id int64) (*models.Message, error)
	UpdateContent(id int64, content string) error
	UpdateStatus(id int64, status string) error
	Delete(id int64) error
	GetConversation(user1ID, user2ID uint, page, perPage int) ([]models.Message, error)
}

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository() MessageRepository {
	return &messageRepository{db: database.DB}
}

func (r *messageRepository) Create(message *models.Message) (int64, error) {
	query := `INSERT INTO messages (sender_id, receiver_id, content, status, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	now := time.Now()
	res, err := stmt.Exec(message.SenderID, message.ReceiverID, message.Content, message.Status, now, now)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *messageRepository) FindByID(id int64) (*models.Message, error) {
	query := `
		SELECT 
			m.id, m.sender_id, m.receiver_id, m.content, m.status, m.created_at, m.updated_at,
			s.id, s.name, s.username,
			r.id, r.name, r.username
		FROM messages m
		JOIN users s ON m.sender_id = s.id
		JOIN users r ON m.receiver_id = r.id
		WHERE m.id = ?`
	row := r.db.QueryRow(query, id)

	var msg models.Message
	err := row.Scan(
		&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Status, &msg.CreatedAt, &msg.UpdatedAt,
		&msg.Sender.ID, &msg.Sender.Name, &msg.Sender.Username,
		&msg.Receiver.ID, &msg.Receiver.Name, &msg.Receiver.Username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) UpdateContent(id int64, content string) error {
	query := `UPDATE messages SET content = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, content, time.Now(), id)
	return err
}

func (r *messageRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE messages SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *messageRepository) Delete(id int64) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *messageRepository) GetConversation(user1ID, user2ID uint, page, perPage int) ([]models.Message, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	query := `
		SELECT 
			m.id, m.sender_id, m.receiver_id, m.content, m.status, m.created_at, m.updated_at,
			s.id, s.name, s.username,
			r.id, r.name, r.username
		FROM messages m
		JOIN users s ON m.sender_id = s.id
		JOIN users r ON m.receiver_id = r.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, user1ID, user2ID, user2ID, user1ID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Status, &msg.CreatedAt, &msg.UpdatedAt,
			&msg.Sender.ID, &msg.Sender.Name, &msg.Sender.Username,
			&msg.Receiver.ID, &msg.Receiver.Name, &msg.Receiver.Username,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}