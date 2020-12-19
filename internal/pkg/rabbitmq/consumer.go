package rabbitmq

import "time"

type ConsumerConfig struct {
	ExchangeName  string
	ExchangeType  string
	RoutingKey    string
	QueueName     string
	ConsumerName  string
	ConsumerCount int
	PrefetchCount int
	Reconnect     struct {
		MaxAttempt int
		Interval   time.Duration
	}
}
