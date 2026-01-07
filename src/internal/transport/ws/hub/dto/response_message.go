package dto

type ResponseMessage struct {
	Type   string `json:"type"`
	ChatID uint   `json:"-"`
	Data   any    `json:"data"`
}
