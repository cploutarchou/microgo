package main

import (
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

	// Routes for static files.
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
