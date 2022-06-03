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
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
func (h *Handlers) GoHome(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
func (h *Handlers) JetHome(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	_data := "Test session data."

	h.APP.Session.Put(r.Context(), "test", _data)
	val := h.APP.Session.GetString(r.Context(), "test")
	vars := make(jet.VarMap)
	vars.Set("test", val)
	err := h.render(w, r, "sessions", vars, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}

func (h *Handlers) Json(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID         int64    `json:"id"`
		Name       string   `json:"name"`
		Activities []string `json:"activities"`
	}
	payload.ID = 100
	payload.Name = "Christos Ploutarchou"
	payload.Activities = []string{"football", "programming", "reading"}

	err := h.APP.WriteJson(w, http.StatusOK, payload)
	if err != nil {
		h.APP.ErrorLog.Println(err)
	}

}

func (h *Handlers) XML(w http.ResponseWriter, r *http.Request) {
	var Payload struct {
		ID         int64    `xml:"id"`
		Name       string   `xml:"name"`
		Activities []string `xml:"activities>activity"`
	}

	var payload Payload
	payload.ID = 100
	payload.Name = "Christos Ploutarchou"
	payload.Activities = []string{"football", "programming", "reading"}

	err := h.APP.WriteXML(w, http.StatusOK, payload)
	if err != nil {
		h.APP.ErrorLog.Println(err)
	}

}

func (h *Handlers) Download(w http.ResponseWriter, r *http.Request) {
	h.APP.SentFile(w, r, "./public/icon", "favicon.png")
}
