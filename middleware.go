// Package logrusmiddleware is a simple net/http middleware for logging
// using logrus
package logrusmiddleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	// Middleware is a middleware handler for HTTP logging
	Middleware struct {
		// Logger is the log.Logger instance used to log messages with the Logger middleware
		Logger *logrus.Logger
		// Name is the name of the application as recorded in latency metrics
		Name string
	}

	// Handler is the actual middleware that handles logging
	Handler struct {
		m         *Middleware
		handler   http.Handler
		component string
	}
)

// Handler create a new handler. component, if set, is emitted in the log messages.
func (m *Middleware) Handler(h http.Handler, component string) *Handler {
	return &Handler{
		m:         m,
		handler:   h,
		component: component,
	}
}

// ServeHTTP calls the "real" handler and logs using the logger
func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()

	lrw := &ResponseWriter{ResponseWriter: rw}
	h.handler.ServeHTTP(lrw, r)

	elapsed := time.Since(start)

	fields := logrus.Fields{
		"status":     lrw.status,
		"method":     r.Method,
		"request":    r.RequestURI,
		"remote":     r.RemoteAddr,
		"elapsed_ms": float64(elapsed) / float64(time.Millisecond),
		"size":       lrw.size,
		"referer":    r.Referer(),
		"user-agent": r.UserAgent(),
	}

	if h.m.Name != "" {
		fields["name"] = h.m.Name
	}

	if h.component != "" {
		fields["component"] = h.component
	}

	if l := h.m.Logger; l != nil {
		l.WithFields(fields).Info("request")
	} else {
		logrus.WithFields(fields).Info("request")
	}
}
