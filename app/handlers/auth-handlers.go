package handlers

import (
	"net/http"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	err := h.APP.Render.Page(w, r, "login", nil, nil)
	if err != nil {
		h.APP.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			h.APP.ErrorLog.Println(err)
			return
		}
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	user, err := h.Models.Users.GetByEmail(email)

	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			h.APP.ErrorLog.Println(err)
			return
		}
		return
	}
	valid, err := user.PasswordMatches(password)
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			h.APP.ErrorLog.Println(err)
			return
		}
		return
	}
	if !valid {
		_, err := w.Write([]byte("Invalid password"))
		if err != nil {
			h.APP.ErrorLog.Println(err.Error())
			return
		}
		return
	}
	h.APP.Session.Put(r.Context(), "userID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.APP.Session.RenewToken(r.Context())
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			h.APP.ErrorLog.Println(err.Error())
			return
		}
		return
	}
	h.APP.Session.Remove(r.Context(), "userID")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}
