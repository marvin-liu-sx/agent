package proxy

import (
	"bee-agent/utils"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"
)

var (
	log = logging.Logger("proxy")
)

func init() {

}

type Manager struct {
	users           map[string]string
	debugAPI        string
	api             string
	port            int64
	serverUrl       string
	serverToken     string
	disableRegister bool
	srv             *http.Server
	quit            chan struct{}
}

func NewManager(api, debugAPI, serverUrl, serverToken string, port int64, disableRegister bool) *Manager {
	m := Manager{
		users:           map[string]string{},
		api:             api,
		debugAPI:        debugAPI,
		serverUrl:       serverUrl,
		serverToken:     serverToken,
		port:            port,
		disableRegister: disableRegister,
		quit:            make(chan struct{}),
	}
	for _, username := range initUsers {
		m.users[username] = utils.RandString(32)
	}
	return &m
}

func (m *Manager) Start() error {

	if m.disableRegister {
		log.Warnf("The agent will not be registered, and the server cannot access the agent!")
	} else {
		if err := utils.Retry(3, time.Second, m.register); err != nil {
			log.Errorf("register err: %s", err)
			return err
		}
	}
	return m.forwardAPI()
}

func (m *Manager) Stop() {
	close(m.quit)
	log.Info("Shutting down api server ...")

	if m.srv == nil {
		return
	}

	if err := m.srv.Shutdown(context.TODO()); err != nil {
		log.Errorf("shutting down api server failed: %s", err)
	}
}

func (m *Manager) forwardAPI() error {

	r := mux.NewRouter()

	r.Use(m.authMiddleware)

	fwd, _ := forward.New()

	api := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// let us forward this request to another server

		prefix := "/api"

		if !strings.HasPrefix(req.URL.Path, prefix) {
			http.NotFound(w, req)
			return
		}

		to := testutils.ParseURI(m.api)
		to.Path = path.Join(to.Path, req.URL.Path[len(prefix):])

		log.Infof("forward api: %s -> %s", req.URL.String(), to.String())

		req.URL = to
		req.RequestURI = req.URL.Path

		fwd.ServeHTTP(w, req)
	})

	debug := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// let us forward this request to another server

		prefix := "/debug"

		if !strings.HasPrefix(req.URL.Path, prefix) {
			http.NotFound(w, req)
			return
		}

		to := testutils.ParseURI(m.debugAPI)
		to.Path = path.Join(to.Path, req.URL.Path[len(prefix):])

		log.Infof("forward debug api: %s -> %s", req.URL.String(), to.String())

		req.URL = to
		req.RequestURI = req.URL.Path

		fwd.ServeHTTP(w, req)
	})

	r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		match, _ := regexp.MatchString("/api/.*", r.URL.Path)
		return match
	}).HandlerFunc(api)

	r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		match, _ := regexp.MatchString("/debug/.*", r.URL.Path)
		return match
	}).HandlerFunc(debug)

	r.HandleFunc("/admin/password", m.getPassword)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", m.port),
		Handler: r,
	}
	m.srv = s

	log.Infof("listen on: %s", s.Addr)
	for k, v := range m.users {
		log.Infow("authorization", "username", k, "pass", v)
	}
	return m.srv.ListenAndServe()
}
