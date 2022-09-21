package handlers

import (
	"test/data"
	microGo "github.com/cploutarchou/MicroGO"
	"net/http"
)

type Handlers struct {
	APP    *microGo.MicroGo
	Models data.Models
}

// Home Request Handler for home package.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	_data := map[string]interface{}{
		"Title": "HomePage",
	}
	err := h.render(w, r, "index", "main", nil, _data)
	if err != nil {
		h.APP.ErrorLog.Println("Unable to render page :", err)
	}
}
