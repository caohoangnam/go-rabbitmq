package http

import (
	"net/http"

	"github.com/working/go-rabbitmq/internal/pkg/rabbitmq"
	"github.com/working/go-rabbitmq/internal/user"
)

type Router struct {
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		http.NewServeMux(),
	}
}

func (r *Router) RegisterUser(rabbitmq *rabbitmq.RabbitMQ) {
	create := user.NewCreate(rabbitmq)

	r.Handler("/users/create", create.Handle)
}
