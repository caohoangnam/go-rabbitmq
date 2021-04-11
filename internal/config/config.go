package config

import (
	"time"

	"github.com/working/go-rabbitmq/internal/pkg/rabbitmq"
	"github.com/working/go-rabbitmq/internal/user"
)

type Config struct {
	HTTPAddress string
	RabbitMQ    rabbitmq.Config
	UserAMQP    user.AMQPConfig
	Consumer    rabbitmq.ConsumerConfig
}

func New() Config {
	var conf Config

	conf.HTTPAddress = ":8080"

	conf.RabbitMQ.Schema = "amqp"
	conf.RabbitMQ.Username = "guest"
	conf.RabbitMQ.Password = "guest"
	conf.RabbitMQ.Host = "0.0.0.0"
	conf.RabbitMQ.Port = "5672"
	conf.RabbitMQ.VHost = "my_app"
	conf.RabbitMQ.ConnectingName = "MY_APP"
	conf.RabbitMQ.ChannelNotifyTimeout = 100 * time.Millisecond
	conf.RabbitMQ.Reconnect.Interval = 500 * time.Millisecond
	conf.RabbitMQ.Reconnect.MaxAttempt = 7200

	// Config user amqp
	conf.UserAMQP.Create.ExchangeName = "user"
	conf.UserAMQP.Create.ExchangeType = "direct"
	conf.UserAMQP.Create.RoutingKey = "create"
	conf.UserAMQP.Create.QueueName = "user_create"

	conf.Consumer.ExchangeName = "user"
	conf.Consumer.ExchangeType = "direct"
	conf.Consumer.RoutingKey = "create"
	conf.Consumer.QueueName = "user_create"
	conf.Consumer.ConsumerCount = 1
	conf.Consumer.PrefetchCount = 1
	conf.Consumer.Reconnect.Interval = 1 * time.Second
	conf.Consumer.Reconnect.MaxAttempt = 60

	return conf
}
