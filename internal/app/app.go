package app

import "net/http"

type App struct {
	server *http.Server
}

func New(server *http.Server) App {
	return App{
		server: server,
	}
}

func (a App) Run() error {
	return a.server.ListenAndServe()
}
