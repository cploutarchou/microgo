package rapiditas

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

//routes Return a Mux object that implements the Router interface.
func (r *Rapiditas) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if r.Debug {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)

	return mux
}
