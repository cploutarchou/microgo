package microGo

import (
	"github.com/justinas/nosurf"
	"net/http"
	"strconv"
	"time"
)

func (m *MicroGo) LoadSessions(next http.Handler) http.Handler {
	m.InfoLog.Println("Triggering LoadSessions for MicroGo.")
	return m.Session.LoadAndSave(next)
}

func (m *MicroGo) NoSurf(next http.Handler) http.Handler {
	csrf := nosurf.New(next)
	isSecure, _ := strconv.ParseBool(m.config.cookie.secure)
	csrf.ExemptGlob("/api/*")
	csrf.SetBaseCookie(http.Cookie{
		Name:       "",
		Value:      "",
		Path:       "",
		Domain:     m.config.cookie.domain,
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     0,
		Secure:     isSecure,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   nil,
	})
	return csrf
}
