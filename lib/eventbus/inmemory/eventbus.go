package inmemory

import (
	"fmt"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"log"
)

func NewEventBus() *EventBus {
	log.Println("starting inmemory eventbus")
	return &EventBus{
		handlers: make(map[EventKey]HandlerGroups),
	}
}

type EventBus struct {
	handlers map[EventKey]HandlerGroups
}

func (e *EventBus) Close() error {
	log.Println("closing inmemory eventbus")
	return nil
}

func (e *EventBus) Publish(event domain.Event) error {
	if handlerGroups := e.handlers[eventKey(event.GetDefinition())]; handlerGroups != nil {
		for _, handler := range handlerGroups.SelectHandlers() {
			if err := handler.GetHandling()(event); err != nil {
				handler.GetErrorHandling()(event, err)
			}
		}
	}

	return nil
}

func (e *EventBus) Listen(handlers ...domain.EventHandler) error {
	for _, handler := range handlers {
		key := eventKey(handler.GetEvent())
		handlerGroups := e.handlers[key]
		if handlerGroups == nil {
			handlerGroups = NewHandlerGroups()
		}

		handlerGroups.AddEventHandler(handler)
		e.handlers[key] = handlerGroups
	}

	return nil
}

type EventKey string

func eventKey(event domain.EventDefinition) EventKey {
	return EventKey(fmt.Sprintf("%s/%s", event.GetDomain(), event.GetName()))
}
