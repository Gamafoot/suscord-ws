package broker

import "suscord_ws/internal/domain/eventbus"

type Broker interface {
	Consume(queue string) error
	SetEventbus(eventbus eventbus.Bus)
	Close() error
}
