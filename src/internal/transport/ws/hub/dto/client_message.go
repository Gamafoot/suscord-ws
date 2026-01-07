package dto

import "encoding/json"

type ClientMessage struct {
	Type   string          `json:"type"`
	ChatID uint            `json:"chat_id,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}
