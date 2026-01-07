package eventbus

import (
	"context"
	"suscord_ws/internal/domain/eventbus"
)

type Eventbus struct {
	handlers map[string]eventbus.Handler
}

func NewEventbus() *Eventbus {
	return &Eventbus{make(map[string]eventbus.Handler)}
}

func (e *Eventbus) Subscribe(name string, handle eventbus.Handler) {
	e.handlers[name] = handle
}

func (e *Eventbus) Publish(ctx context.Context, name string, payload eventbus.Payload) error {
	return e.handlers[name](ctx, payload)
}
