package main

import (
	"log"
	"os"

	"github.com/cploutarchou/rapiditas"
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

	rap.AppName = "myapp"
	rap.InfoLog.Println("DEBUG mode is set to : ", rap.Debug)

	app := &application{
		App: rap,
	}

	return app
}
