package message

import (
	"suscord_ws/internal/domain/broker/message/model"
)

type MessageCreated struct {
	model.Message
}
