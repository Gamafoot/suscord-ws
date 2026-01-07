package message

import (
	"suscord_ws/internal/domain/broker/message/model"
)

type ChatUpdated struct {
	model.Chat
}
