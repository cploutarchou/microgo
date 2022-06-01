package middleware

import (
	"app/data"
	microGo "cloud0.christosploutarchou.com/cploutarchou/MicroGO"
)

type Middleware struct {
	App    *microGo.MicroGo
	Models data.Models
}
