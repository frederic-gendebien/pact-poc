package rabbitmq

import (
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"github.com/streadway/amqp"
)

func createTopic(channel *amqp.Channel, spec domain.EventDefinition) error {
	return channel.ExchangeDeclare(
		spec.GetDomain(), //name
		"topic",          //kind
		true,             //durable
		false,            //auto_delete
		false,            //internal
		false,
		nil,
	)
}
