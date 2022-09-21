package main

import (
	"test/data"
	"test/handlers"
	"test/middleware"
	microGo "github.com/cploutarchou/MicroGO"
	"log"
	"os"
)

func initApplication() *test {
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

	rap.AppName = "test"

	testMiddleware := &middleware.Middleware{
		App: rap,
	}
	testHandlers := &handlers.Handlers{
		APP: rap,
	}

	test := &test{
		App:        rap,
		Handlers:   testHandlers,
		Middleware: testMiddleware,
	}
	test.App.Routes = test.routes()
	test.Models = data.New(test.App.DB.Pool)
	testHandlers.Models = test.Models
	test.Middleware.Models = test.Models

	return test
}
