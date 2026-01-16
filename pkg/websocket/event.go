package websocket

import "github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"

const (
	EventMessageSent    = "message_sent"
	EventMessageUpdated = "message_updated"
	EventMessageDeleted = "message_deleted"
)

type Event struct {
	Type    string         `json:"type"`
	Payload models.Message `json:"payload"`
}
