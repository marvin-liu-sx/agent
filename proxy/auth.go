package proxy

import (
	auth "github.com/abbot/go-http-auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

const (
	adminUsername = "admin"
	apiUsername   = "api"
	debugUsername = "debug"
)

var initUsers = []string{adminUsername, apiUsername, debugUsername}

func (m *Manager) authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sp := strings.Split(r.URL.Path, "/")
		if len(sp) < 2 {
			http.NotFound(w, r)
			return
		}
		authenticator := auth.NewBasicAuthenticator(sp[1], m.secret)

		if u := authenticator.CheckAuth(r); u == "" {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Manager) secret(user, realm string) string {

	if pass, ok := m.users[user]; ok {

		hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return ""
		}
		if user == adminUsername {
			return string(hash)
		}
		if realm == user {
			return string(hash)
		}
		return ""
	}
	return ""
}
