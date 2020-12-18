package user

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/working/go-rabbitmq/internal/pkg/rabbitmq"
)

type AMQPConfig struct {
	Create struct {
		ExchangeName string
		ExchangeType string
		RoutingKey   string
		QueueName    string
	}
}

type AMQP struct {
	config   AMQPConfig
	rabbitmq *rabbitmq.RabbitMQ
}

func NewAMQP(config AMQPConfig, rabbitmq *rabbitmq.RabbitMQ) AMQP {
	return AMQP{
		config:   config,
		rabbitmq: rabbitmq,
	}
}

func (a *AMQP) Setup() (err error) {
	channel, err := a.rabbitmq.Channel()
	if err != nil {
		err = errors.Wrap(err, "Failed to open channel")
		return
	}
	defer channel.Close()

	if err := a.declareCreate(channel); err != nil {
		return
	}
	return
}

func (a *AMQP) declareCreate(channel *amqp.Channel) (err error) {
	// ExchangeDeclare declares an exchange on the server.
	// If the exchange does not already exist, the server will create it.
	// If the exchange exists, the server verifies that it is of the provided type, durability and auto-delete flags.
	err = channel.ExchangeDeclare(
		a.config.Create.ExchangeName,
		a.config.Create.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		err = errors.Wrap(err, "Failed to declare exchange")
		return
	}

	_, err = channel.QueueDeclare(
		a.config.Create.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-queue-model": "lazy"},
	)
	if err != nil {
		err = errors.Wrap(err, "Failed to declare queue")
		return
	}

	err = channel.QueueBind(
		a.config.Create.QueueName,
		a.config.Create.RoutingKey,
		a.config.Create.ExchangeName,
		false,
		false,
	)
	if err != nil {
		err = errors.Wrap(err, "Failed to bind queue")
		return
	}
	return
}
