package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

//routes The application routes.
func (a *application) routes() *chi.Mux {
	// Note : All the middlewares must come before the routes

	// Add routes here
	a.App.Routes.Get("/", a.Handlers.Home)
	a.App.Routes.Get("/go-page", a.Handlers.GoHome)
	a.App.Routes.Get("/jet-page", a.Handlers.JetHome)
	a.App.Routes.Get("/sessions", a.Handlers.SessionTest)
	a.App.Routes.Get("/testdb", func(w http.ResponseWriter, r *http.Request) {
		query := "select id,first_name from user where id = 1"
		rows, _ := a.App.DB.Pool.QueryContext(r.Context(), query)
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			a.App.ErrorLog.Print(err)
			return
		}
		fmt.Print(w, "%d %s", id, name)
	})

	// Routes for static files.
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
