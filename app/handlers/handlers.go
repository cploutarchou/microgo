package handlers

import (
	"github.com/cploutarchou/microGo"
	"net/http"
)

type Handlers struct {
	APP *microGo.MicroGo
}

// Home Request Handler for home package.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.APP.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
