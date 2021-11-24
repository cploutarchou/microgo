package session

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"reflect"
	"testing"
)

func TestSessionInitSession(t *testing.T) {

	c := &Session{
		CookieLifetime: "100",
		CookiePersist:  "true",
		CookieName:     "microGo",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}
	var sm *scs.SessionManager
	ses := c.InitializeSession()
	var sessionKind reflect.Kind
	var sessionType reflect.Type
	rv := reflect.ValueOf(ses)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Print("In Loop : ", rv.Kind(), rv.Type(), rv)
		sessionKind = rv.Kind()
		sessionType = rv.Type()
		rv.Type()
		rv = rv.Elem()

	}
	if !rv.IsValid() {
		t.Error("Invalid type or Kind! Kind:", rv.Kind(), "type:", rv.Type())
	}
	if sessionKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind returned testing cookie session. Expected", reflect.ValueOf(sm).Kind(), "and got", sessionKind)
	}

	if sessionType != reflect.ValueOf(sm).Type() {
		t.Error("wrong type returned testing cookie session. Expected", reflect.ValueOf(sm).Type(), "and got", sessionType)
	}
}
