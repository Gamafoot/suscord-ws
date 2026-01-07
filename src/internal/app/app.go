package app

import (
	"context"
	"log"
	"suscord_ws/internal/config"
	"suscord_ws/internal/domain/broker"
	"suscord_ws/internal/infrastructure/eventbus"
	"time"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
	errwg "golang.org/x/sync/errgroup"
)

type App struct {
	cfg      *config.Config
	echo     *echo.Echo
	broker   broker.Broker
	wsServer *websocketServer
}

func NewApp() (*App, error) {
	cfg := config.GetConfig()

	storage, err := NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	echo := echo.New()
	eventbus := eventbus.NewEventbus()

	wsServer := NewWebsocketServer(
		cfg,
		echo,
		storage,
		eventbus,
	)

	broker, err := NewBrokerConsumer(cfg.Broker.Addr, eventbus)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:      cfg,
		echo:     echo,
		broker:   broker,
		wsServer: wsServer,
	}, nil
}

func (app *App) Run() error {
	wg, _ := errwg.WithContext(context.Background())

	wg.Go(func() error {
		port := app.cfg.Server.Port

		if err := app.wsServer.Run(port); err != nil {
			return err
		}

		return nil
	})

	wg.Go(func() error {
		log.Println("RabbitMQ consumer is running")
		if err := app.broker.Consume("chat.ws"); err != nil {
			log.Printf("broker err: %+v\n", err)
			return err
		}
		return nil
	})

	return wg.Wait()
}

func (app *App) Stop() error {
	ctx, cansel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cansel()

	wg, ctx := errwg.WithContext(ctx)

	wg.Go(func() error {
		return app.broker.Close()
	})

	wg.Go(func() error {
		if err := app.wsServer.echo.Shutdown(ctx); err != nil {
			return pkgErrors.WithStack(err)
		}
		return nil
	})

	return wg.Wait()
}
