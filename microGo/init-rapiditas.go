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
	cel := &rapiditas.Rapiditas{}
	err = cel.New(path)
	if err != nil {
		log.Fatal(err)
	}

	cel.AppName = "myapp"
	cel.Debug = true

	app := &application{
		App: cel,
	}

	return app
}
