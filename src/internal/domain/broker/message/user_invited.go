package message

type UserInvited struct {
	ChatID uint   `json:"chat_id"`
	UserID uint   `json:"user_id"`
	Code   string `json:"code"`
}
