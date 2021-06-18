package proxy

import (
	"encoding/json"
	"net/http"
)

func (m *Manager) getPassword(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if pass, ok := m.users[username]; ok {
		s := struct {
			Password string `json:"password"`
		}{
			Password: pass,
		}
		v, _ := json.Marshal(s)
		w.Write(v)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}
	return
}
