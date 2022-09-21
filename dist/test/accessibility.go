package main

import (
	"net/http"
)

//get helper function for test get request.
func (a *test) get(s string, h http.HandlerFunc) {
	a.App.Routes.Get(s, h)
}

//get helper function for test post request.
func (a *test) post(s string, h http.HandlerFunc) {
	a.App.Routes.Post(s, h)

}

//get helper function for middlewares
func (a *test) use(m ...func(handler http.Handler) http.Handler) {
	a.App.Routes.Use(m...)
}
