package middleware

import (
	"app/data"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (m *Middleware) CheckRemember(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if user is not logged in
		if !m.App.Session.Exists(r.Context(), "userID") {
			cookie, err := r.Cookie(fmt.Sprintf("_%s_remember", m.App.AppName))
			if err != nil {
				n.ServeHTTP(w, r)
			} else {
				//cookie has already exist

				key := cookie.Value
				var u data.User
				if len(key) > 0 {
					// validate cookie value
					split := strings.Split(key, "|")
					uid, hash := split[0], split[1]
					id, _ := strconv.Atoi(uid)
					validHash := u.CheckForRememberToken(id, hash)
					if !validHash {
						m.deleteRememberCookie(w, r)
						m.App.Session.Put(r.Context(), "error", "You have logged out")
						n.ServeHTTP(w, r)
					} else {
						//logged in user
						user, _ := u.Get(id)
						m.App.Session.Put(r.Context(), "userID", user.ID)
						m.App.Session.Put(r.Context(), "remember_token", hash)
						n.ServeHTTP(w, r)
					}
				} else {
					// leftover cookie (Browser session is still open)
					m.deleteRememberCookie(w, r)
					n.ServeHTTP(w, r)
				}
			}

		} else {
			n.ServeHTTP(w, r)
		}
	})
}

func (m *Middleware) deleteRememberCookie(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	//delete the cookie
	newCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_remember", m.App.AppName),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-50 * time.Hour),
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   m.App.Session.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &newCookie)
	// logout the user

	m.App.Session.Remove(r.Context(), "userID")
	m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

}
