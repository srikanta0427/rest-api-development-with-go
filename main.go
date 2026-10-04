package main

import (
	"log"

	"github.com/srikanta0427/rest_api_design/app"
)

func main() {
	cfg := app.NewConfig() // set the server to listen on port
	appl := app.NewApplication(cfg)
	errRun := appl.Run()
	if errRun != nil {
		log.Fatal(errRun)
	}
}
