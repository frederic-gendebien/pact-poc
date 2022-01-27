package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"github.com/streadway/amqp"
	"log"
)

const (
	url = "RABBITMQ_URL"
)

func NewEventBus(configuration config.Configuration) *EventBus {
	log.Println("connecting to rabbitmq")
	url := configuration.GetStringOrCrash(url)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("could not connect to rabbitmq: %v", err)
	}

	return &EventBus{
		connection: conn,
	}
}

type EventBus struct {
	connection *amqp.Connection
}

func (e *EventBus) Close() error {
	log.Println("closing rabbitmq eventbus")
	return e.connection.Close()
}

func (e *EventBus) Publish(ctx context.Context, event domain.Event) error {
	payload, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}

	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	if err := createTopic(channel, event.GetDefinition()); err != nil {
		return err
	}

	return channel.Publish(
		event.GetDefinition().GetDomain(), //exchange
		event.GetDefinition().GetName(),   //key
		true,                              //mandatory
		false,                             //immediate
		amqp.Publishing{
			Headers:     newHeaders(event.GetDefinition()),
			ContentType: "application/json",
			Body:        payload,
		},
	)
}

func (e *EventBus) Listen(ctx context.Context, listenerName string, eventHandlers ...domain.EventHandler) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	queueName, err := createQueue(channel, listenerName, eventHandlers...)
	if err != nil {
		return err
	}

	messages, err := consumeQueue(channel, queueName)
	if err != nil {
		return err
	}

	handlers := handlerMap(eventHandlers...)
	for message := range messages {
		processMessage(message, handlers)
	}

	return fmt.Errorf("no more message available")
}

func handlerMap(handlers ...domain.EventHandler) map[string]domain.EventHandler {
	m := make(map[string]domain.EventHandler)
	for _, eventHandler := range handlers {
		m[eventHandler.GetEventDefinition().GetName()] = eventHandler
	}

	return m
}

func processMessage(message amqp.Delivery, handlers map[string]domain.EventHandler) {
	messageDomain := message.Headers[EventDomain].(string)
	messageType := message.Headers[EventType].(string)
	if handler, present := handlers[messageType]; present {
		event := handler.GetEventDefinition().GetType()
		if err := json.Unmarshal(message.Body, event); err != nil {
			log.Printf("could not unmarshal event: %v", err)
			handler.HandleError(message.Body, err)
		}

		if err := handler.ProcessEvent(event); err != nil {
			log.Printf("could not process message from domain (%s) of type (%s): %v", messageDomain, messageType, err)
			handler.HandleError(event, err)
		}
	} else {
		log.Printf("skip message from domain (%s) of type (%s)", messageDomain, messageType)
	}

	_ = message.Ack(false)
}
