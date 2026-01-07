package hub

import (
	"suscord_ws/internal/domain/entity"
	"suscord_ws/internal/transport/ws/hub/dto"
)

func (hub *Hub) joinChatRoom(chatID uint, client *HubClient) error {
	if client == nil {
		return nil
	}

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	if _, exists := hub.chatRooms[chatID]; !exists {
		hub.chatRooms[chatID] = make(map[uint]bool)
	}

	hub.chatRooms[chatID][client.ID] = true
	client.rooms[chatID] = true

	return nil
}

func (hub *Hub) leaveChatRoom(roomID, clientID uint) {
	hub.mutex.Lock()

	if room, exists := hub.chatRooms[roomID]; exists {
		delete(room, clientID)
		if len(room) == 0 {
			delete(hub.chatRooms, roomID)
		}
	}

	if client, exists := hub.clients[clientID]; exists {
		delete(client.rooms, roomID)
	}

	hub.mutex.Unlock()

	hub.broadcastToChatRoomExcept(roomID, clientID, &dto.ResponseMessage{
		Type: "user_left",
		Data: map[string]interface{}{"user_id": clientID},
	})
}

func (hub *Hub) deleteChatRoom(roomID uint) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.chatRooms[roomID] = make(map[uint]bool)

	for clientID := range hub.clients {
		delete(hub.clients[clientID].rooms, roomID)
	}
}

func (hub *Hub) joinToUserChatRooms(client *HubClient, chats []*entity.Chat) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for _, chat := range chats {
		if _, exists := hub.chatRooms[chat.ID]; !exists {
			hub.chatRooms[chat.ID] = make(map[uint]bool)
		}

		hub.chatRooms[chat.ID][client.ID] = true
		client.rooms[chat.ID] = true
	}

	return nil
}

func (hub *Hub) broadcastToChatRoom(roomID uint, message interface{}) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.chatRooms[roomID]; exists {
		for userID := range room {
			if client, exists := hub.clients[userID]; exists {
				client.SendMessage(message)
			}
		}
	}
}

func (hub *Hub) broadcastToChatRoomExcept(roomID, exceptUserID uint, message any) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.chatRooms[roomID]; exists {
		for userID := range room {
			if userID != exceptUserID {
				if client, exists := hub.clients[userID]; exists {
					client.SendMessage(message)
				}
			}
		}
	}
}

func (hub *Hub) broadcastToSFURoom(roomID uint, message any) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.sfuRooms[roomID]; exists {
		for userID := range room {
			if client, exists := hub.clients[userID]; exists {
				client.SendMessage(message)
			}
		}
	}
}

func (hub *Hub) broadcastToSFURoomExcept(roomID, exceptUserID uint, message any) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.sfuRooms[roomID]; exists {
		for userID := range room {
			if userID != exceptUserID {
				if client, exists := hub.clients[userID]; exists {
					client.SendMessage(message)
				}
			}
		}
	}
}
