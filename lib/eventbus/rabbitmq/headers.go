package rabbitmq

import (
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"github.com/streadway/amqp"
)

const (
	EventDomain = "event.domain"
	EventType   = "event.type"
)

func newHeaders(event domain.EventDefinition) amqp.Table {
	m := make(map[string]interface{})
	m[EventDomain] = event.GetDomain()
	m[EventType] = event.GetName()
	return m
}
