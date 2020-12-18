package http

import (
	"net/http"
)

func NewServer(address string, router *Router) *http.Server {
	return &http.Server{
		Addr:    address,
		Handler: router,
	}
}
