package hub

import (
	"context"
	"suscord_ws/internal/domain/broker/message"
	"suscord_ws/internal/domain/broker/queue"
	"suscord_ws/internal/domain/eventbus"
	"suscord_ws/internal/transport/ws/hub/dto"
)

func (hub *Hub) registerEvents(bus eventbus.Bus) {
	bus.Subscribe(queue.MessageCreated, hub.onMessageCreated)
	bus.Subscribe(queue.MessageUpdated, hub.onMessageUpdated)
	bus.Subscribe(queue.MessageDeleted, hub.onMessageDeleted)
	bus.Subscribe(queue.ChatUpdated, hub.onChatUpdated)
	bus.Subscribe(queue.ChatDeleted, hub.onDeleteChat)
	bus.Subscribe(queue.UserInvited, hub.onUserInvited)
	bus.Subscribe(queue.UserJoinedGroupChat, hub.onUserJoinedGroupChat)
	bus.Subscribe(queue.UserJoinedPrivateChat, hub.onUserJoinedPrivateChat)
	bus.Subscribe(queue.UserLeft, hub.onUserLeft)
}

func (hub *Hub) onMessageCreated(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.MessageCreated](payload)
	if err != nil {
		return err
	}

	hub.broadcastToChatRoom(data.ChatID, &dto.ResponseMessage{
		Type:   "message",
		ChatID: data.ChatID,
		Data:   data,
	})

	return nil
}

func (hub *Hub) onMessageUpdated(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.MessageUpdated](payload)
	if err != nil {
		return err
	}

	hub.broadcastToChatRoomExcept(data.ChatID, data.UserID, &dto.ResponseMessage{
		Type:   "message_update",
		ChatID: data.ChatID,
		Data:   data,
	})

	return nil
}

func (hub *Hub) onMessageDeleted(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.MessageDeleted](payload)
	if err != nil {
		return err
	}

	hub.broadcastToChatRoomExcept(data.ChatID, data.ExceptUserID, &dto.ResponseMessage{
		Type:   "message_delete",
		ChatID: data.ChatID,
		Data:   data,
	})

	return nil
}

func (hub *Hub) onUserInvited(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.UserInvited](payload)
	if err != nil {
		return err
	}

	if client, exists := hub.clients[data.UserID]; exists {
		client.SendMessage(&dto.ResponseMessage{
			Type: "invite_to_chat",
			Data: map[string]string{
				"code": data.Code,
			},
		})
	}

	return nil
}

func (hub *Hub) onChatUpdated(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.ChatUpdated](payload)
	if err != nil {
		return err
	}

	hub.broadcastToChatRoom(data.Chat.ID, &dto.ResponseMessage{
		Type: "update_group_chat",
		Data: data,
	})

	return nil
}

func (hub *Hub) onUserJoinedGroupChat(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.UserJoinedGroupChat](payload)
	if err != nil {
		return err
	}

	if client, exists := hub.clients[data.User.ID]; exists {
		hub.joinChatRoom(data.ChatID, client)
		hub.broadcastToChatRoomExcept(data.ChatID, data.User.ID, &dto.ResponseMessage{
			Type:   "new_user_in_chat",
			ChatID: data.ChatID,
			Data:   data,
		})
	}

	return nil
}

func (hub *Hub) onUserJoinedPrivateChat(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.UserJoinedPrivateChat](payload)
	if err != nil {
		return err
	}

	if client, exists := hub.clients[data.User.ID]; exists {
		hub.joinChatRoom(data.ChatID, client)
		client.SendMessage(&dto.ResponseMessage{
			Type: "joined_chat",
			Data: data,
		})
	}

	return nil
}

func (hub *Hub) onUserLeft(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.UserLeft](payload)
	if err != nil {
		return err
	}

	hub.leaveChatRoom(data.ChatID, data.UserID)

	hub.broadcastToChatRoomExcept(data.ChatID, data.UserID, &dto.ResponseMessage{
		Type:   "user_left",
		ChatID: data.ChatID,
		Data:   data,
	})

	return nil
}

func (hub *Hub) onDeleteChat(ctx context.Context, payload eventbus.Payload) error {
	data, err := unmarshal[message.ChatDeleted](payload)
	if err != nil {
		return err
	}

	hub.broadcastToChatRoom(data.ID, &dto.ResponseMessage{
		Type: "delete_chat",
		Data: data,
	})
	hub.deleteChatRoom(data.ID)

	return nil
}
