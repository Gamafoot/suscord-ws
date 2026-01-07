package app

import (
	"suscord_ws/internal/domain/broker"
	"suscord_ws/internal/domain/eventbus"
	"suscord_ws/internal/infrastructure/broker/rabbitmq"
)

func NewBrokerConsumer(addr string, eventbus eventbus.Bus) (broker.Broker, error) {
	broker, err := rabbitmq.NewBroker(addr)
	if err != nil {
		return nil, err
	}

	broker.SetEventbus(eventbus)

	return broker, nil
}
