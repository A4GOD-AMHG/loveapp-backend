package models

import (
	"time"
)

type Message struct {
	ID         uint      `json:"id"`
	SenderID   uint      `json:"sender_id"`
	Sender     User      `json:"sender"`
	ReceiverID uint      `json:"receiver_id"`
	Receiver   User      `json:"receiver"`
	Content    string    `json:"content"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
