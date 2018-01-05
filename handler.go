package logrusmiddleware

import (
	"net/http"
)

// ResponseWriter exposes the number of
type ResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

// Write is a wrapper for the "real" ResponseWriter.Write
func (w *ResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		w.status = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// WriteHeader is a wrapper around ResponseWriter.WriteHeader
func (w *ResponseWriter) WriteHeader(s int) {
	w.status = s
	w.ResponseWriter.WriteHeader(s)
}
