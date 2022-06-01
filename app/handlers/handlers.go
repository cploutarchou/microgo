package handlers

import (
	"app/data"
	microGo "cloud0.christosploutarchou.com/cploutarchou/MicroGO"
	"github.com/CloudyKit/jet/v6"
	"net/http"
)

type Handlers struct {
	APP    *microGo.MicroGo
	Models data.Models
}

// Home Request Handler for home package.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.APP.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
func (h *Handlers) GoHome(w http.ResponseWriter, r *http.Request) {
	err := h.APP.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
func (h *Handlers) JetHome(w http.ResponseWriter, r *http.Request) {
	err := h.APP.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	data := "Test session data."

	h.APP.Session.Put(r.Context(), "test", data)
	val := h.APP.Session.GetString(r.Context(), "test")
	vars := make(jet.VarMap)
	vars.Set("test", val)
	err := h.APP.Render.Page(w, r, "sessions", vars, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
