package main

import (
	"log"

	"github.com/working/go-rabbitmq/internal/app"
	"github.com/working/go-rabbitmq/internal/config"
	"github.com/working/go-rabbitmq/internal/pkg/http"
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

	// AMQP Setup
	userAMQP := user.NewAMQP(conf.UserAMQP, rbt)
	if err := userAMQP.Setup(); err != nil {
		log.Fatalln(err)
	}

	// HTTP Router
	router := http.NewRouter()
	router.RegisterUser()

	// HTTP Server
	server := http.NewServer(conf.HTTPAddress, router)

	// Run
	err := app.New(server).Run()
	if err != nil {
		log.Fatalln(err)
	}
}
