package main

import (
	"app/data"
	"app/handlers"
	"github.com/cploutarchou/microGo"
)

type application struct {
	App      *microGo.MicroGo
	Handlers *handlers.Handlers
	Models   data.Models
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
