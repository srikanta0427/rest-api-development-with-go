package app

import (
	"fmt"
	"net/http"

	"github.com/srikanta0427/rest_api_design/config"
	"github.com/srikanta0427/rest_api_design/controller"
	db "github.com/srikanta0427/rest_api_design/db/repos"
	"github.com/srikanta0427/rest_api_design/router"
	"github.com/srikanta0427/rest_api_design/services"
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
	Store  db.Storage
}

func NewApplication(cfg Config) *Application {
	return &Application{
		Config: cfg,
		Store:  *db.NewStorage(),
	}
}

func (app *Application) Run() error {
	ur := db.NewUserRepository()
	us := services.NewUserService(ur)
	uc := controller.NewUserController(us)
	ud := router.NewUserRouter(uc)
	server := http.Server{
		Addr:    app.Config.Addr,
		Handler: router.SetUpRouter(ud),
	}
	fmt.Println("Server is listening on " + app.Config.Addr)
	return server.ListenAndServe()
}
