package middleware

import (
	"app/data"
	"github.com/cploutarchou/microGo"
)

type Middleware struct {
	App    *microGo.MicroGo
	Models data.Models
}
