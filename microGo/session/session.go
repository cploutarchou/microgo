package session

import (
	"github.com/alexedwards/scs/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	IsCookieSecure string
}

func (s *Session) InitializeSession() *scs.SessionManager {
	var persist, secure bool
	// Session lifetime
	minutes, err := strconv.Atoi(s.CookieLifetime)
	if err != nil {
		minutes = 60
	}
	// Should cookies persists?
	if strings.ToLower(s.CookiePersist) == "true" {
		persist = true
	} else {
		persist = false
	}

	// is cookie secure?
	if strings.ToLower(s.IsCookieSecure) == "true" {
		secure = true
	}

	// initiate session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = s.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = s.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	// select session store

	switch strings.ToLower(s.SessionType) {
	case "redis":
	case "mysql", "mariadb":
	case "postgres", "postgresql":
	default:

	}
	return session
}
