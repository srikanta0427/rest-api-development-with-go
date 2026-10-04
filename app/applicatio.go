package app

import (
	"fmt"
	"net/http"

	"github.com/srikanta0427/rest_api_design/config"
)

type Config struct {
	Addr string
}

// Constructor for type Config struct

func NewConfig() Config {
	port := config.GetString("ADDR", ":8080")
	return Config{
		Addr: port,
	}
}

type Application struct {
	Config Config
}

func NewApplication(cfg Config) *Application {
	return &Application{
		Config: cfg,
	}
}

func (app *Application) Run() error {
	server := http.Server{
		Addr:    app.Config.Addr,
		Handler: nil,
	}
	fmt.Println("Server is listening on " + app.Config.Addr)
	return server.ListenAndServe()
}
