package message

type UserLeft struct {
	ChatID uint `json:"chat_id"`
	UserID uint `json:"user_id"`
}
