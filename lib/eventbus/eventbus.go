package eventbus

import (
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/inmemory"
	"io"
	"log"
)

const (
	Mode         = "EVENTBUS_MODE"
	ModeInMemory = "inmemory"
)

func NewEventBus(configuration config.Configuration) EventBus {
	mode := configuration.GetMandatoryValue(Mode)
	switch mode {
	case ModeInMemory:
		return inmemory.NewEventBus()
	default:
		log.Fatalf("unknown eventbus mode: %s", mode)
		return nil
	}
}

type EventBus interface {
	io.Closer
	Publish(domain.Event) error
	Listen(handlers ...domain.EventHandler) error
}
