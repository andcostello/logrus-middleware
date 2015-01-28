package logrusmiddleware

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
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
		http.ResponseWriter
		status    int
		size      int
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

// Write is a wrapper for the "real" Write
func (h *Handler) Write(b []byte) (int, error) {
	if h.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		h.status = http.StatusOK
	}
	size, err := h.ResponseWriter.Write(b)
	h.size += size
	return size, err
}

func (h *Handler) WriteHeader(s int) {
	h.ResponseWriter.WriteHeader(s)
	h.status = s
}

func (h *Handler) Header() http.Header {
	return h.ResponseWriter.Header()
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()

	h.handler.ServeHTTP(rw, r)

	latency := time.Since(start)

	fields := logrus.Fields{
		"status":   h.status,
		"method":   r.Method,
		"request":  r.RequestURI,
		"remote":   r.RemoteAddr,
		"duration": latency.Seconds(),
		"size":     h.size,
	}

	if h.m.Name != "" {
		fields["name"] = h.m.Name
	}

	if h.component != "" {
		fields["component"] = h.component
	}

	h.m.Logger.WithFields(fields).Info("completed handling request")
}
