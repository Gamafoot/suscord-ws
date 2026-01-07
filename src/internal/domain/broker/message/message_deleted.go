package message

type MessageDeleted struct {
	ChatID       uint
	MessageID    uint
	ExceptUserID uint
}
