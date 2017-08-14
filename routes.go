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

func Router(m *web.Mux, cn *c.Config, l *logrus.Logger, s func() map[interface{}]interface{}) {
	// Handle /config route
	m.Get("/config", func(c web.C, w http.ResponseWriter, r *http.Request) {
		http.Error(w, config.Dump(cn), http.StatusOK)
	})

	// Handle /debug routes
	m.Post("/debug", func(c web.C, w http.ResponseWriter, r *http.Request) {
		if l.Level != logrus.DebugLevel {
			l.Level = logrus.DebugLevel
			l.Debug("Debug level set")

			http.Error(w, "ok", http.StatusOK)
		} else {
			http.Error(w, "error", http.StatusOK)
		}

		http.Error(w, config.Dump(cn), http.StatusOK)
	})
	m.Delete("/debug", func(c web.C, w http.ResponseWriter, r *http.Request) {
		if l.Level == logrus.DebugLevel {
			l.Level = logrus.InfoLevel
			l.Info("Debug level unset")

			http.Error(w, "ok", http.StatusOK)
		} else {
			http.Error(w, "error", http.StatusOK)
		}

		http.Error(w, config.Dump(cn), http.StatusOK)
	})

	// Handle /status route
	m.Handle("/status", func(c web.C, w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(s())
		if err != nil {
			l.Warning(fmt.Sprintf("Error marshaling status to json: %v", err))
		} else {
			http.Error(w, string(b), http.StatusOK)
		}
	})
}
