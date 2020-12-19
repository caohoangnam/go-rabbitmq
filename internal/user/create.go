package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/working/go-rabbitmq/internal/pkg/rabbitmq"
)

type Create struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewCreate(rabbitmq *rabbitmq.RabbitMQ) Create {
	return Create{
		rabbitmq: rabbitmq,
	}
}

func (c Create) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("ID")
	if err := c.publish(id); err != nil {
		log.Println(errors.Wrap(err, fmt.Sprintf("Failed to create %s", id)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("created", id)
	w.WriteHeader(http.StatusAccepted)
}

func (c Create) publish(msg string) error {
	channel, err := c.rabbitmq.Channel()
	if err != nil {
		return errors.Wrap(err, "Failed to open channel")
	}
	defer channel.Close()

	if err := channel.Confirm(false); err != nil {
		return errors.Wrap(err, "Failed to put channel is confirm mode")
	}

	err = channel.Publish(
		"user",
		"create",
		true,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			MessageId:    "A-UNIQUE-ID",
			ContentType:  "text/plain",
			Body:         []byte(msg),
		},
	)
	if err != nil {
		return errors.Wrap(err, "Failed to publish message")
	}

	select {
	case ntf := <-channel.NotifyPublish(make(chan amqp.Confirmation, 1)):
		if !ntf.Ack {
			return errors.New("Failed to deliver msg to exchange/queue")
		}
	case <-channel.NotifyReturn(make(chan amqp.Return)):
		return errors.New("Failed to delivery msg to exchange/queue")
		//	case <-time.After(c.rabbitmq.ChannelNotifyTimeout):
		//		return errors.New("Failed msg delivery confirm to exchange/queue timed out")
	}

	return nil
}
