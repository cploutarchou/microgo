package middleware

import (
	"test/data"
	microGo "github.com/cploutarchou/MicroGO"
)

type Middleware struct {
	App    *microGo.MicroGo
	Models data.Models
}
