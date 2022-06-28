package handlers

import (
	"app/data"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
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
	// Check if user set remember me flag and
	if r.Form.Get("remember") == "remember" {
		randomStr := h.createRandomString(12)
		hasher := sha256.New()
		_, err := hasher.Write([]byte(randomStr))
		if err != nil {
			h.APP.ErrorStatus(w, http.StatusBadRequest)
			return
		}
		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		rm := data.RememberToken{}
		err = rm.InsertToken(user.ID, sha)
		if err != nil {
			h.APP.ErrorStatus(w, http.StatusBadRequest)
			return
		}

		// set cookie to remember
		expire := time.Now().Add(365 * 24 * 60 * 60 * time.Second)
		cookie := http.Cookie{
			Name:     fmt.Sprintf("_%s_remember", h.APP.AppName),
			Value:    fmt.Sprintf("%d|%s", user.ID, sha),
			Path:     "/",
			HttpOnly: true,
			Expires:  expire,
			Domain:   h.APP.Session.Cookie.Domain,
			MaxAge:   31535000,
			Secure:   h.APP.Session.Cookie.Secure,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		h.APP.Session.Put(r.Context(), "remember_token", sha)

	}
	h.APP.Session.Put(r.Context(), "userID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {

	// delete remember token  from session if exists
	if h.APP.Session.Exists(r.Context(), "remember_token") {
		rt := data.RememberToken{}
		_ = rt.Delete(h.APP.Session.GetString(r.Context(), "remember_token"))
	}
	cookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_remember", h.APP.AppName),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-100 * time.Hour),
		Domain:   h.APP.Session.Cookie.Domain,
		MaxAge:   -1,
		Secure:   h.APP.Session.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	h.APP.Session.RenewToken(r.Context())
	h.APP.Session.Remove(r.Context(), "userID")
	h.APP.Session.Remove(r.Context(), "remember_token")
	h.APP.Session.Destroy(r.Context())
	h.APP.Session.RenewToken(r.Context())
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
