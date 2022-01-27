package eventbus

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/usecase"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/events"
	eventbus "github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"log"
)

const (
	ListenerName = "projection"
)

func NewUserRegisteredHandler(useCase usecase.UserProjectionUseCase) eventbus.EventHandler {
	return NewListener(
		events.NewUserRegistered{},
		func(event interface{}) error {
			return useCase.AddUser(context.Background(), projectionUser(event.(*events.NewUserRegistered).User))
		},
		logError(),
	)
}

func UserDeletedHandler(useCase usecase.UserProjectionUseCase) eventbus.EventHandler {
	return NewListener(
		events.UserDeleted{},
		func(event interface{}) error {
			return useCase.DeleteUserById(context.Background(), model.UserId(event.(*events.UserDeleted).UserId))
		},
		logError(),
	)
}

func logError() func(event interface{}, err error) {
	return func(event interface{}, err error) {
		log.Printf("error processing event: %v: %v", event, err)
	}
}

func NewListener(
	event eventbus.EventDefinition,
	handling func(interface{}) error,
	errorHandling func(interface{}, error),
) *Listener {
	return &Listener{
		event:         event,
		handling:      handling,
		errorHandling: errorHandling,
	}
}

type Listener struct {
	event         eventbus.EventDefinition
	handling      func(interface{}) error
	errorHandling func(interface{}, error)
}

func (l *Listener) GetName() string {
	return "projection"
}

func (l *Listener) GetEventDefinition() eventbus.EventDefinition {
	return l.event
}

func (l *Listener) ProcessEvent(event interface{}) error {
	return l.handling(event)
}

func (l *Listener) HandleError(event interface{}, err error) {
	l.errorHandling(event, err)
}
