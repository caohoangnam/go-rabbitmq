package main

import (
	"log"

	"github.com/working/go-rabbitmq/internal/config"
	"github.com/working/go-rabbitmq/internal/pkg/rabbitmq"
	"github.com/working/go-rabbitmq/internal/user"
)

func main() {

	//Load Config
	conf := config.New()

	// Load RabbitMQ
	rbt := rabbitmq.New(conf.RabbitMQ)
	if err := rbt.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer rbt.Shutdown()

	consumer := user.NewConsumer(conf.Consumer, rbt)
	if err := consumer.Start(); err != nil {
		log.Fatalln(err)
	}
}
