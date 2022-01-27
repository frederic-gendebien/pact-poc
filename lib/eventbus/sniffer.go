package eventbus

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"log"
	"sync"
)

func NewEventSniffer(eventBus EventBus) *EventSniffer {
	return &EventSniffer{
		lock:     &sync.RWMutex{},
		eventBus: eventBus,
	}
}

type EventSniffer struct {
	lock     *sync.RWMutex
	eventBus EventBus
	events   []interface{}
}

func (e *EventSniffer) Listen(eventDefinition domain.EventDefinition) error {
	return e.eventBus.Listen(context.Background(),
		"event-sniffer",
		NewEventListener(e, eventDefinition),
	)
}

func (e *EventSniffer) Clear() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.events = nil
}

func (e *EventSniffer) AddEvent(event interface{}) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.events = append(e.events, event)
}

func (e *EventSniffer) GetEvents() []interface{} {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.events
}

func NewEventListener(eventSniffer *EventSniffer, eventDefinition domain.EventDefinition) *EventListener {
	return &EventListener{
		eventSniffer:    eventSniffer,
		eventDefinition: eventDefinition,
	}
}

type EventListener struct {
	eventSniffer    *EventSniffer
	eventDefinition domain.EventDefinition
}

func (e EventListener) GetEventDefinition() domain.EventDefinition {
	return e.eventDefinition
}

func (e EventListener) ProcessEvent(event interface{}) error {
	e.eventSniffer.AddEvent(event)

	return nil
}

func (e EventListener) HandleError(event interface{}, err error) {
	log.Printf("could not process event: %v: %v", event, err)
}
