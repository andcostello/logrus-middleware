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
func (w *ResponseWriter) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.size += n
	return
}

// WriteHeader is a wrapper around ResponseWriter.WriteHeader
func (w *ResponseWriter) WriteHeader(s int) {
	w.status = s
	w.ResponseWriter.WriteHeader(s)
}
