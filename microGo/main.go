package main

import (
	"github.com/cploutarchou/rapiditas"
	"microGo/handlers"
)

type application struct {
	App      *rapiditas.Rapiditas
	Handlers *handlers.Handlers
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
