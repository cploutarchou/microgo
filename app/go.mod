module app

go 1.17

replace github.com/cploutarchou/microGo => ./../microGo

require (
	github.com/cploutarchou/microGo v0.0.0
	github.com/go-chi/chi/v5 v5.0.7
)

require github.com/joho/godotenv v1.4.0 // indirect
