package handlers

import (
	"github.com/justinas/nosurf"
	"net/http"
)

func (h *Handlers) ShowCachePage(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "cache", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) SaveInCache(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		CSRF  string `json:"csrf_token"`
	}

	err := h.APP.ReadJson(w, r, &userInput)
	if err != nil {
		h.APP.Error500(w, r)
		return
	}
	if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
		h.APP.ErrorUnprocessable(w, r)
		return
	}
	err = h.APP.Cache.Set(userInput.Name, userInput.Value)
	if err != nil {
		h.APP.Error500(w, r)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "Saved in cache"

	_ = h.APP.WriteJson(w, http.StatusCreated, resp)
}

func (h *Handlers) GetFromCache(w http.ResponseWriter, r *http.Request) {
	var msg string
	var inCache = true

	var userInput struct {
		Name string `json:"name"`
		CSRF string `json:"csrf_token"`
	}

	err := h.APP.ReadJson(w, r, &userInput)
	if err != nil {
		h.APP.Error500(w, r)
		return
	}
	if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
		h.APP.ErrorUnprocessable(w, r)
		return
	}
	fromCache, err := h.APP.Cache.Get(userInput.Name)
	if err != nil {
		msg = "Not found in cache!"
		inCache = false
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Value   string `json:"value"`
	}

	if inCache {
		resp.Error = false
		resp.Message = "Success"
		resp.Value = fromCache.(string)
	} else {
		resp.Error = true
		resp.Message = msg
	}
	_ = h.APP.WriteJson(w, http.StatusCreated, resp)
}

func (h *Handlers) DeleteFromCache(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Name string `json:"name"`
		CSRF string `json:"csrf_token"`
	}

	err := h.APP.ReadJson(w, r, &userInput)
	if err != nil {
		h.APP.Error500(w, r)
		return
	}
	if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
		h.APP.ErrorUnprocessable(w, r)
		return
	}
	err = h.APP.Cache.Delete(userInput.Name)
	if err != nil {
		h.APP.Error500(w, r)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	resp.Message = "Deleted from cache (if it existed)"

	_ = h.APP.WriteJson(w, http.StatusCreated, resp)
}

func (h *Handlers) EmptyCache(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		CSRF string `json:"csrf_token"`
	}

	err := h.APP.ReadJson(w, r, &userInput)
	if err != nil {
		h.APP.Error500(w, r)
		return
	}
	if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
		h.APP.ErrorUnprocessable(w, r)
		return
	}
	err = h.APP.Cache.Clean()
	if err != nil {
		h.APP.Error500(w, r)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	resp.Message = "Emptied cache!"

	_ = h.APP.WriteJson(w, http.StatusCreated, resp)

}
