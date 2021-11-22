package main

import (
	"github.com/cploutarchou/rapiditas"
	"log"
	"microGo/handlers"
	"os"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init rapiditas
	rap := &rapiditas.Rapiditas{}
	err = rap.New(path)
	if err != nil {
		log.Fatal(err)
	}

	rap.AppName = "app"
	add_handlers := &handlers.Handlers{
		APP: rap,
	}

	app := &application{
		App:      rap,
		Handlers: add_handlers,
	}
	app.App.Routes = app.routes()

	return app
}
