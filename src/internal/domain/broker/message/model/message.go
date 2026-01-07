package model

import (
	"time"
)

type Message struct {
	ID          uint          `json:"id"`
	ChatID      uint          `json:"chat_id"`
	UserID      uint          `json:"user_id"`
	Content     string        `json:"content"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Attachments []*Attachment `json:"attachments"`
}
