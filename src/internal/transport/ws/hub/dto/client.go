package dto

import (
	"suscord_ws/internal/transport/ws/hub/model"
	"suscord_ws/pkg/urlpath"
)

type Client struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

func NewClient(client *model.Client, mediaURL string) *Client {
	if client == nil {
		return nil
	}

	return &Client{
		ID:        client.ID,
		Username:  client.Username,
		AvatarUrl: urlpath.GetMediaURL(mediaURL, client.AvatarPath),
	}
}
