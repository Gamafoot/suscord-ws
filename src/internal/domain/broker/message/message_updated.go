package message

import (
	"suscord_ws/internal/domain/broker/message/model"
)

type MessageUpdated struct {
	model.Message
}
