package main

import (
	"app/handlers"
	"github.com/cploutarchou/microGo"
)

type application struct {
	App      *microGo.MicroGo
	Handlers *handlers.Handlers
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
