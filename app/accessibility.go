package main

import (
	"net/http"
)

//get helper function for app get request.
func (a *application) get(s string, h http.HandlerFunc) {
	a.App.Routes.Get(s, h)
}

//get helper function for app post request.
func (a *application) post(s string, h http.HandlerFunc) {
	a.App.Routes.Get(s, h)

}

//get helper function for middlewares
func (a *application) use(m ...func(handler http.Handler) http.Handler) {
	a.App.Routes.Use(m...)
}
