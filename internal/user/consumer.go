package user

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"github.com/working/go-rabbitmq/internal/pkg/rabbitmq"
)

type UserConsumer struct {
	config   rabbitmq.ConsumerConfig
	rabbitmq *rabbitmq.RabbitMQ
}

func NewConsumer(config rabbitmq.ConsumerConfig, rabbitmq *rabbitmq.RabbitMQ) *UserConsumer {
	return &UserConsumer{
		config:   config,
		rabbitmq: rabbitmq,
	}
}

func (c *UserConsumer) Start() error {
	conn, err := c.rabbitmq.Connection()
	if err != nil {
		return err
	}
	defer conn.Close()

	go c.closedConnectionListener(conn.NotifyClose(make(chan *amqp.Error)))

	chn, err := conn.Channel()
	if err != nil {
		return err
	}

	// Exchange declare
	err = chn.ExchangeDeclare(
		c.config.ExchangeName, // name
		c.config.ExchangeType, // kind
		true,                  // durable
		false,                 // autoDelete
		false,                 // internal
		false,                 // noAwait
		nil,                   // Table
	)
	if err != nil {
		return err
	}

	// Queue declare
	_, err = chn.QueueDeclare(
		c.config.QueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   //amqp.Table{"x-queue-mode": "lazy"},
	)
	if err != nil {
		return err
	}

	// Queue Bind
	err = chn.QueueBind(
		c.config.QueueName,
		c.config.RoutingKey,
		c.config.ExchangeName,
		false, //noWait
		nil,   //Table
	)
	if err != nil {
		return err
	}

	//http://www.rabbitmq.com/blog/2012/04/25/rabbitmq-performance-measurements-part-2/
	err = chn.Qos(c.config.PrefetchCount, 0, false)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	for i := 1; i <= c.config.ConsumerCount; i++ {
		id := i
		go c.consumer(chn, id)
	}
	<-forever
	return nil
}

func (c *UserConsumer) closedConnectionListener(closed <-chan *amqp.Error) {
	log.Println("INFO: Watching close connection")

	err := <-closed
	if err != nil {
		log.Println("INFO: Closed connection ->", err.Error())

		for i := 0; i < c.config.Reconnect.MaxAttempt; i++ {
			log.Println("INFO: Attempt to reconnect")

			time.Sleep(c.config.Reconnect.Interval)
			err := c.rabbitmq.Connect()
			if err != nil {
				continue
			}

			log.Println("INFO: Reconnected")
			err = c.Start()
			if err != nil {
				continue
			}
			break
		}
	} else {
		log.Println("INFO: Connection closed normally, will not reconnect")
		os.Exit(0)
	}
}

func (c *UserConsumer) consumer(channel *amqp.Channel, id int) {
	msgs, err := channel.Consume(
		c.config.QueueName,
		"", //	fmt.Sprintf("%s (%d/%d)", c.config.ConsumerName, id, c.config.ConsumerCount),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(fmt.Sprintf("CRITICAL: Unable to start consumer (%d/%d)", id, c.config.ConsumerCount))

		return
	}

	log.Println("[", id, "] Running ...")
	log.Println("[", id, "] Press CTRL+C to exit ...")

	for msg := range msgs {
		log.Println("[", id, "] Consumed:", string(msg.Body))

		if err := msg.Ack(false); err != nil {
			// TODO: Should DLX the message
			// http://www.inanzzz.com/index.php/post/1p7m/creating-a-rabbitmq-dlx-dead-letter-exchange-example-with-golang
			log.Println("unable to acknowledge the message, dropped", err)
		}
	}

	log.Println("[", id, "] Exiting ...")
}
