package rabbitmq

import "encoding/json"

type MessageBody struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
