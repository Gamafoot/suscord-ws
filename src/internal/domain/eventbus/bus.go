package eventbus

import "context"

type Payload []byte

type Handler func(ctx context.Context, payload Payload) error

type Bus interface {
	Subscribe(name string, handle Handler)
	Publish(ctx context.Context, name string, payload Payload) error
}
