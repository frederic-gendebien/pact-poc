package eventbus

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/inmemory"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/rabbitmq"
	"io"
	"log"
)

const (
	mode         = "EVENTBUS_MODE"
	modeInMemory = "inmemory"
	modeRabbitMQ = "rabbitmq"
)

func NewEventBus(configuration config.Configuration) EventBus {
	mode := configuration.GetStringOrCrash(mode)
	switch mode {
	case modeInMemory:
		return inmemory.NewEventBus()
	case modeRabbitMQ:
		return rabbitmq.NewEventBus(configuration)
	default:
		log.Fatalf("unknown eventbus mode: %s", mode)
		return nil
	}
}

type EventBus interface {
	io.Closer
	Publish(ctx context.Context, event domain.Event) error
	Listen(ctx context.Context, listenerName string, eventHandlers ...domain.EventHandler) error
}
