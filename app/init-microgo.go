package main

import (
	"app/data"
	"app/handlers"
	"app/middleware"
	microGo "cloud0.christosploutarchou.com/cploutarchou/MicroGO"
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

	appMiddleware := &middleware.Middleware{
		App: rap,
	}
	appHandlers := &handlers.Handlers{
		APP: rap,
	}

	app := &application{
		App:        rap,
		Handlers:   appHandlers,
		Middleware: appMiddleware,
	}
	app.App.Routes = app.routes()
	app.Models = data.New(app.App.DB.Pool)
	appHandlers.Models = app.Models
	app.Middleware.Models = app.Models

	return app
}
