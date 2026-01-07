package app

import (
	"context"
	"fmt"
	"suscord_ws/internal/config"
	"suscord_ws/internal/domain/eventbus"
	"suscord_ws/internal/domain/storage"
	"suscord_ws/internal/transport/ws"
	"suscord_ws/internal/transport/ws/hub"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
	errwg "golang.org/x/sync/errgroup"
)

type websocketServer struct {
	cfg  *config.Config
	echo *echo.Echo
	hub  *hub.Hub
}

func NewWebsocketServer(
	cfg *config.Config,
	echo *echo.Echo,
	storage storage.Storage,
	eventbus eventbus.Bus,
) *websocketServer {
	hubInstance := hub.NewHub(cfg, storage, eventbus)

	server := &websocketServer{
		cfg:  cfg,
		echo: echo,
		hub:  hubInstance,
	}

	handler := ws.NewHandler(hubInstance)
	handler.InitRoutes(server.echo)

	return server
}

func (s *websocketServer) Run(port int) error {
	wg, _ := errwg.WithContext(context.Background())

	wg.Go(func() error {
		s.hub.Run()
		return nil
	})

	wg.Go(func() error {
		if err := s.echo.Start(fmt.Sprintf(":%d", port)); err != nil {
			return pkgErrors.WithStack(err)
		}
		return nil
	})

	return wg.Wait()
}
