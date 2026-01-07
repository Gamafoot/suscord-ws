package message

import "suscord_ws/internal/domain/broker/message/model"

type UserJoinedPrivateChat struct {
	ChatID uint       `json:"chat_id"`
	User   model.User `json:"user"`
}
