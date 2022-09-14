package session

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/uptrace/bun"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	IsCookieSecure string
	DBPool         *bun.DB
	RedisPool      *redis.Pool
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
		session.Store = redisstore.New(s.RedisPool)
	case "mysql", "mariadb":
		session.Store = mysqlstore.New(s.DBPool.DB)
	case "postgres", "postgresql":
		session.Store = postgresstore.New(s.DBPool.DB)
	default:

	}
	return session
}
