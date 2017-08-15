package go_monitoring_routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/endeveit/go-snippets/config"
	c "github.com/robfig/config"
	"github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

type Router struct {
	mux        *web.Mux
	conf       *c.Config
	log        *logrus.Logger
	version    interface{}
	statusFunc func() map[string]interface{}
	healthFunc func() bool
}

func NewRouter() *Router {
	r := Router{}

	return &r
}

// Set mux for router
func (r *Router) Mux(m *web.Mux) *Router {
	r.mux = m

	return r
}

// Set /version route handler
func (r *Router) Version(v interface{}) *Router {
	r.version = v

	r.mux.Get("/version", func(c web.C, w http.ResponseWriter, req *http.Request) {
		http.Error(w, fmt.Sprintf("%v", r.version), http.StatusOK)
	})

	return r
}

// Set /config route handler
func (r *Router) Config(cn *c.Config) *Router {
	r.conf = cn

	r.mux.Get("/config", func(c web.C, w http.ResponseWriter, req *http.Request) {
		http.Error(w, config.Dump(r.conf), http.StatusOK)
	})

	return r
}

// Set /debug route handler
func (r *Router) Logger(l *logrus.Logger) *Router {
	r.log = l

	r.mux.Post("/debug", func(c web.C, w http.ResponseWriter, req *http.Request) {
		if r.log.Level != logrus.DebugLevel {
			r.log.Level = logrus.DebugLevel
			r.log.Debug("Debug level set")

			http.Error(w, "ok", http.StatusOK)
		} else {
			http.Error(w, "error", http.StatusOK)
		}
	})
	r.mux.Delete("/debug", func(c web.C, w http.ResponseWriter, req *http.Request) {
		if r.log.Level == logrus.DebugLevel {
			r.log.Level = logrus.InfoLevel
			r.log.Info("Debug level unset")

			http.Error(w, "ok", http.StatusOK)
		} else {
			http.Error(w, "error", http.StatusOK)
		}
	})

	return r
}

// Set /status route handler
func (r *Router) StatusFunc(f func() map[string]interface{}) *Router {
	r.statusFunc = f

	r.mux.Handle("/status", func(c web.C, w http.ResponseWriter, req *http.Request) {
		b, err := json.Marshal(r.statusFunc())
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling status to json: %v", err), http.StatusOK)
		} else {
			http.Error(w, string(b), http.StatusOK)
		}
	})

	return r
}

// Set /health route handler
func (r *Router) HealthFunc(f func() bool) *Router {
	r.healthFunc = f

	r.mux.Handle("/health", func(c web.C, w http.ResponseWriter, req *http.Request) {
		if r.healthFunc() {
			http.Error(w, "green", http.StatusOK)
		} else {
			http.Error(w, "red", http.StatusServiceUnavailable)
		}
	})

	return r
}
