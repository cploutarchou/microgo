package main

import (
	"test/data"
	"test/handlers"
	"test/middleware"
	microGo "github.com/cploutarchou/MicroGO"
)

type test struct {
	App        *microGo.MicroGo
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
