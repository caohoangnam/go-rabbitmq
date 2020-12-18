package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Config struct {
	Schema               string
	Username             string
	Password             string
	Host                 string
	Port                 string
	VHost                string
	ConnectingName       string
	ChannelNotifyTimeout time.Duration
	Reconnect            struct {
		Interval   time.Duration
		MaxAttempt int
	}
}

type RabbitMQ struct {
	// A RWMutex is a reader/writer mutual exclusion lock.
	// The lock can be held by an arbitrary number of readers or a single writer.
	// The zero value for a RWMutex is an unlocked mutex.
	// A RWMutex must not be copied after first use.
	mux                  sync.RWMutex
	config               Config
	dialConfig           amqp.Config
	connection           *amqp.Connection
	ChannelNotifyTimeout time.Duration
}

func New(config Config) *RabbitMQ {
	return &RabbitMQ{
		config: config,
		dialConfig: amqp.Config{
			Properties: amqp.Table{
				"connection_name": config.ConnectingName,
			},
		},
	}

}

func (r *RabbitMQ) Connect() (err error) {
	conn, err := amqp.DialConfig(fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		r.config.Schema,
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
		r.config.VHost,
	), r.dialConfig)
	if err != nil {
		return
	}

	r.connection = conn

	go r.reconnect()

	return
}

func (r *RabbitMQ) Connection() *amqp.Connection { return r.connection }

func (r *RabbitMQ) Shutdown() error {
	if r.connection != nil {
		return r.connection.Close()
	}
}

func (r *RabbitMQ) Channel() (*amqp.Channel, error) {
	if r.connection == nil {
		if err := r.Connect(); err != nil {
			return nil, errors.New("Connection isn't open")
		}
	}

	channel, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *RabbitMQ) reconnect() {
WATCH:
	connErr := <-r.connection.NotifyClose(make(chan *amqp.Error))
	if connErr != nil {
		log.Println("CRITICAL: Connection dropped, reconnecting")

		var err error

		for i = 1; i <= r.config.Reconnect.MaxAttempt; i++ {
			r.mux.RLock()
			r.connection, err = amqp.DialConfig(fmt.Sprintf(
				"%s://%s:%s@%s:%s/%s",
				r.config.Schema,
				r.config.Username,
				r.config.Password,
				r.config.Host,
				r.config.Port,
				r.config.Vhost,
			), r.dialConfig)
			r.mux.RUnlock()

			if err != nil {
				log.Println("INFO: Reconnection")

				goto WATCH
			}
			time.Sleep(r.config.Reconnect.Interval)
		}
		log.Println(errors.Wrap(err, "CRITICAL: Failed to reconnect"))
	} else {
		log.Println("INFO: Connection dropped normally, will not reconnection")
	}
}
