package hub

import (
	"suscord_ws/internal/transport/ws/hub/model"

	"github.com/gorilla/websocket"
)

type HubClient struct {
	*model.Client
	conn  *websocket.Conn
	rooms map[uint]bool
}

func NewHubClient(client *model.Client, conn *websocket.Conn) *HubClient {
	return &HubClient{
		Client: client,
		conn:   conn,
		rooms:  make(map[uint]bool),
	}
}

func (c *HubClient) SendMessage(messageData interface{}) error {
	if c.conn != nil {
		return c.conn.WriteJSON(messageData)
	}
	return nil
}
