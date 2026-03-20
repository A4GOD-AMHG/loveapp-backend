// Paquete repository implementa el acceso a datos de la aplicación,
// encapsulando todas las operaciones de base de datos.
package repository

import (
	"database/sql"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

// MessageRepository define el contrato de operaciones sobre mensajes en la base de datos.
type MessageRepository interface {
	// Create inserta un nuevo mensaje y retorna su ID generado.
	Create(message *models.Message) (int64, error)
	// FindByID busca un mensaje por su ID junto con los datos del remitente y destinatario.
	FindByID(id int64) (*models.Message, error)
	// UpdateContent actualiza el contenido textual de un mensaje existente.
	UpdateContent(id int64, content string) error
	// UpdateStatus actualiza el estado de entrega/lectura de un mensaje.
	UpdateStatus(id int64, status string) error
	// Delete elimina un mensaje de la base de datos por su ID.
	Delete(id int64) error
	// GetConversation retorna el historial de mensajes paginado entre dos usuarios.
	GetConversation(user1ID, user2ID uint, page, perPage int) ([]models.Message, error)
	// CountUnreadByReceiver retorna la cantidad de mensajes no leídos de un usuario.
	CountUnreadByReceiver(userID uint) (int, error)
}

// messageRepository es la implementación concreta de MessageRepository usando SQLite.
type messageRepository struct {
	db *sql.DB // Conexión a la base de datos
}

// NewMessageRepository crea y retorna una nueva instancia de messageRepository
// conectada a la base de datos global de la aplicación.
func NewMessageRepository() MessageRepository {
	return &messageRepository{db: database.DB}
}

// Create inserta un nuevo mensaje en la base de datos con su estado inicial "sent".
// Retorna el ID generado por la base de datos.
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

// FindByID busca un mensaje por su ID, haciendo JOIN con la tabla de usuarios
// para incluir los datos del remitente y el destinatario.
// Retorna nil si el mensaje no existe.
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

// UpdateContent actualiza el contenido textual de un mensaje y su timestamp de actualización.
func (r *messageRepository) UpdateContent(id int64, content string) error {
	query := `UPDATE messages SET content = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, content, time.Now(), id)
	return err
}

// UpdateStatus actualiza el estado del mensaje (ej. "sent" → "delivered" → "read")
// y su timestamp de actualización.
func (r *messageRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE messages SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// Delete elimina permanentemente un mensaje de la base de datos por su ID.
func (r *messageRepository) Delete(id int64) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// GetConversation retorna el historial de mensajes entre dos usuarios con paginación.
// Los mensajes se ordenan por fecha de creación descendente (más recientes primero).
// Si page o perPage son menores a 1, se usan valores predeterminados (1 y 10 respectivamente).
func (r *messageRepository) GetConversation(user1ID, user2ID uint, page, perPage int) ([]models.Message, error) {
	// Aplicar valores predeterminados de paginación
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	// Obtener mensajes en ambas direcciones entre los dos usuarios
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

	// Mapear cada fila al modelo Message con datos de remitente y destinatario
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

// CountUnreadByReceiver retorna la cantidad de mensajes cuyo destinatario es userID
// y que todavía no han llegado al estado "read".
func (r *messageRepository) CountUnreadByReceiver(userID uint) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM messages
		WHERE receiver_id = ? AND status != 'read'
	`, userID).Scan(&count)

	return count, err
}
