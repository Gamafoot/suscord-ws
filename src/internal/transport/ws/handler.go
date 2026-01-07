package ws

import (
	"suscord_ws/internal/transport/middleware"
	"suscord_ws/internal/transport/ws/hub"

	"github.com/labstack/echo/v4"
)

type handler struct {
	hub *hub.Hub
}

func NewHandler(hub *hub.Hub) *handler {
	return &handler{
		hub: hub,
	}
}

func (h *handler) InitRoutes(route *echo.Echo) {
	route.GET("/ws", h.hub.ServeWS, middleware.Logger)
}
