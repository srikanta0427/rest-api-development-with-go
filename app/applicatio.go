package app

import (
	"fmt"
	"net/http"

	"github.com/srikanta0427/rest_api_design/config"
	config2 "github.com/srikanta0427/rest_api_design/config/db"
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
	Config  Config
	Storage db.Storage
}

func NewApplication(cfg Config) *Application {

	return &Application{
		Config:  cfg,
		Storage: *db.NewStorage(),
	}
}

func (app *Application) Run() error {
	dbconfig, err := config2.SetUpDB()
	if err != nil {
		return err
	}

	// User Repository
	ur := db.NewUserRepository(dbconfig)
	us := services.NewUserService(ur)
	uc := controller.NewUserController(us)
	ud := router.NewUserRouter(uc)

	// email
	emailService := services.NewEmailService()

	// Organizer
	organizationRepo := db.NewOrganizer(dbconfig)
	organizerService := services.NewOrganizerService(organizationRepo, emailService)
	organizerController := controller.NewOrganizerController(organizerService)
	organizerRouter := router.NewOrganizerRouter(organizerController)


	httpRouter := router.SetUpRouter(ud, organizerRouter)

	// server
	server := http.Server{
		Addr:    app.Config.Addr,
		Handler: httpRouter,
	}

	fmt.Println("Server is listening on " + app.Config.Addr)

	// listening
	return server.ListenAndServe()
}
