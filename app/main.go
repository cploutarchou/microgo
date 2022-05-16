package main

import (
	"app/data"
	"app/handlers"
	"app/middleware"
	"github.com/cploutarchou/microGo"
)

type application struct {
	App        *microGo.MicroGo
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
