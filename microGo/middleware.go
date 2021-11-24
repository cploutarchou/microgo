package microGo

import "net/http"

func (m *MicroGo) LoadSessions(next http.Handler) http.Handler {
	m.InfoLog.Println("Triggering LoadSessions for MicroGo.")
	return m.Session.LoadAndSave(next)
}
