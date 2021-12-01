package main

import (
	"app/data"
	"app/handlers"
	"github.com/cploutarchou/microGo"
	"log"
	"os"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init microGo
	rap := &microGo.MicroGo{}
	err = rap.New(path)
	if err != nil {
		log.Fatal(err)
	}

	rap.AppName = "app"
	addHandlers := &handlers.Handlers{
		APP: rap,
	}

	app := &application{
		App:      rap,
		Handlers: addHandlers,
	}
	app.App.Routes = app.routes()
	app.Models = data.New(app.App.DB.Pool)

	return app
}
