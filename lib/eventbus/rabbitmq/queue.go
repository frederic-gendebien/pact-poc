package rabbitmq

import (
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"github.com/streadway/amqp"
)

func createQueue(channel *amqp.Channel, listenerName string, eventHandlers ...domain.EventHandler) (string, error) {
	queueName := listenerName
	if err := declareQueue(channel, queueName); err != nil {
		return "", err
	}

	for _, eventHandler := range eventHandlers {
		if err := createTopic(channel, eventHandler.GetEventDefinition()); err != nil {
			deleteQueue(channel, queueName)
			return "", err
		}
		if err := bindQueue(channel, queueName, eventHandler.GetEventDefinition()); err != nil {
			deleteQueue(channel, queueName)
			return "", err
		}
	}

	return queueName, nil
}

func declareQueue(channel *amqp.Channel, queueName string) error {
	_, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

func bindQueue(channel *amqp.Channel, queueName string, spec domain.EventDefinition) error {
	return channel.QueueBind(
		queueName,
		spec.GetName(),   //key
		spec.GetDomain(), // exchange
		true,
		nil,
	)
}

func deleteQueue(channel *amqp.Channel, queueName string) {
	_, _ = channel.QueueDelete(
		queueName,
		false,
		false,
		true,
	)
}

func consumeQueue(channel *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	return channel.Consume(
		queueName,
		"",
		false,
		false,
		true,
		false,
		nil,
	)
}
